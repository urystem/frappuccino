package dal

import (
	"database/sql"
	"errors"
	"slices"

	"frappuccino/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type dalOrder struct {
	database *sqlx.DB
}

type OrderDalInter interface {
	SelectAllOrders() ([]models.Order, error)
	SelectOrder(uint64) (*models.Order, error)
	DeleteOrder(uint64) error
	InsertOrder(*models.Order, *[]models.InventoryUpdate) error
	UpdateOrder(*models.Order) error
	CloseOrder(uint64) error
	SelectAllStatusHistory() ([]models.StatusHistory, error)
}

func ReturnDulOrderDB(db *sqlx.DB) OrderDalInter {
	return &dalOrder{database: db}
}

func (db *dalOrder) SelectAllOrders() ([]models.Order, error) {
	tx, err := db.database.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orders []models.Order

	err = tx.Select(&orders, `SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareNamed(`SELECT product_id, quantity FROM order_items WHERE order_id = :id`)
	if err != nil {
		return nil, nil
	}

	defer stmt.Close()

	for i, v := range orders {
		err = stmt.Select(&orders[i].Items, v)
		if err != nil {
			return nil, err
		}
	}
	return orders, tx.Commit()
}

func (db *dalOrder) SelectOrder(id uint64) (*models.Order, error) {
	tx, err := db.database.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var order models.Order

	err = tx.Get(&order, `SELECT * FROM orders WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	err = tx.Select(&order.Items, `SELECT product_id, quantity FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &order, tx.Commit()
}

func (db *dalOrder) DeleteOrder(id uint64) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	// егер әлі жабылмаған тапсырыс болса inventory ді түгендейді
	status, err := db.getStatus(tx, id)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	if status == "processing" {
		err = db.inventoryRejector(tx, id)
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}
	// жәй өшіре саламыз. order_items тен өзі өшіп кетеді
	_, err = tx.Exec(`DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

func (db *dalOrder) InsertOrder(ord *models.Order, invUpdates *[]models.InventoryUpdate) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		return err
	}

	if err = tx.QueryRow(`
	INSERT INTO orders (customer_name, allergens)
	VALUES($1,$2)
	RETURNING id`, ord.CustomerName, ord.Allergens).Scan(&ord.ID); err != nil {
		return err
	}
	err = db.detectorAndInserterOrderItems(tx, ord, invUpdates)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *dalOrder) UpdateOrder(ord *models.Order) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	status, err := db.getStatus(tx, ord.ID)
	if err != nil {
		return err
	}

	if status != "processing" {
		return errors.New("it is closed order")
	}
	err = db.inventoryRejector(tx, ord.ID)
	if err != nil {
		return err
	}

	// тазалау
	_, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, ord.ID)
	if err != nil {
		return err
	}
	err = db.detectorAndInserterOrderItems(tx, ord, nil)
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(`
	UPDATE orders 
		SET 
			customer_name = :customer_name, 
			allergens = :allergens,
			updated_at = CURRENT_TIMESTAMP
		WHERE id=:id`, ord)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *dalOrder) CloseOrder(id uint64) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	status, err := db.getStatus(tx, id)
	if err != nil {
		return err
	}
	if status != "processing" {
		return models.ErrOrderStatusClosed
	}
	_, err = tx.Exec(`UPDATE orders
		SET status = 'accepted',
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *dalOrder) SelectAllStatusHistory() ([]models.StatusHistory, error) {
	var statusHistory []models.StatusHistory
	err := db.database.Select(&statusHistory, "SELECT * FROM order_status_history ORDER BY updated_at ASC")
	if err != nil {
		return nil, err
	}
	return statusHistory, nil
}

func (db *dalOrder) getStatus(tx *sqlx.Tx, id uint64) (string, error) {
	var status string
	err := tx.Get(&status, `SELECT status FROM orders WHERE id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", models.ErrNotFound
	} else if err != nil {
		return "", err
	}
	return status, nil
}

func (db *dalOrder) inventoryRejector(tx *sqlx.Tx, orderID uint64) error {
	_, err := tx.Exec(`
	UPDATE inventory AS inv
	SET quantity = inv.quantity + (ings.quantity * ord.quantity)
	FROM menu_item_ingredients AS ings
	JOIN order_items AS ord ON ings.product_id = ord.product_id
	WHERE inv.id = ings.inventory_id AND ord.order_id = $1`, orderID)
	// UPDATE inventory AS inv SET quantity = inv.quantity+(i.quantity * o.quantity) FROM menu_item_ingredients AS i JOIN order_items AS o ON i.product_id = o.product_id WHERE o.order_id = 1;
	if err != nil {
		return err
	}

	const queryTransaction string = `
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		SELECT 
			inv.id,
			(ings.quantity * ord.quantity) AS quantity_change,
			'cancelled'::reason_of_inventory_transaction
		FROM 
			inventory inv
		JOIN 
			menu_item_ingredients ings ON inv.id = ings.inventory_id
		JOIN 
			order_items ord ON ings.product_id = ord.product_id
		WHERE 
			ord.order_id = $1`
	_, err = tx.Exec(queryTransaction, orderID)
	return err
}

func (db *dalOrder) mergerInv(in []models.InventoryUpdate, out *[]models.InventoryUpdate) {
	for _, invent := range in {
		var isHere bool
		for i, inv := range *out {
			if inv.InventoryID == invent.InventoryID {
				(*out)[i].QuantityUsed += invent.QuantityUsed
				isHere = true
				break
			}
		}
		if !isHere {
			*out = append(*out, invent)
		}
	}
}

func (db *dalOrder) checkAllergens(orderAller pq.StringArray, menuAller *pq.StringArray) {
	var invalids uint64
	for _, menuAll := range *menuAller {
		if slices.Contains(orderAller, menuAll) {
			(*menuAller)[invalids] = menuAll
			invalids++
		}
	}
	*menuAller = (*menuAller)[:invalids]
}

func (db *dalOrder) detectorAndInserterOrderItems(tx *sqlx.Tx, ord *models.Order, invsUpdatesOriginal *[]models.InventoryUpdate) error {
	// проверяет существует ли в меню через select allergens
	stmt, err := tx.Preparex(`SELECT allergens FROM menu_items WHERE id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// проверяет достаточно ли ингридентов
	const notEnoughInventsQ string = `
	SELECT id, name, ABS(garbage) AS not_enough
		FROM (
  			SELECT 
				inv.id,
				inv.name,
				inv.quantity - ings.quantity * $2 AS garbage
  			FROM 
				inventory inv
  			JOIN 
    			menu_item_ingredients ings ON inv.id = ings.inventory_id
  			WHERE 
    			ings.product_id = $1
		) sub
	WHERE 
		garbage < 0`
	// SELECT id AS inventory_id, ABS(garbage) AS not_enough FROM (SELECT inv.id, inv.quantity - ings.quantity * $2 AS garbage FROM inventory inv JOIN menu_item_ingredients ings ON inv.id = ings.inventory_id WHERE ings.product_id = $1) sub WHERE garbage < 0;

	stmt2, err := tx.Preparex(notEnoughInventsQ)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	// ВСтавляет запись на order_items (Если до этого все items существует и ингридиенты достаточно)
	stmt3, err := tx.Prepare(`INSERT INTO order_items VALUES($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt3.Close()

	// получить отчеть ингридиентов на каждый меню items
	const remaining string = `
	SELECT 
		inv.id,
		inv.name,
		ings.quantity * $2 AS quantity_used,
		inv.quantity - ings.quantity * $2 AS remaining
	FROM 
		inventory inv
	JOIN 
		menu_item_ingredients ings ON inv.id = ings.inventory_id
	WHERE 
		ings.product_id = $1`

	stmt4, err := tx.Preparex(remaining)
	if err != nil {
		return err
	}
	defer stmt4.Close()

	// minus inventory for every menu item
	const minusInvCycle string = `
	UPDATE inventory AS inv
	SET quantity = inv.quantity - (mi.quantity * $2)
	FROM menu_item_ingredients AS mi
	JOIN menu_items AS m ON mi.product_id = m.id
	WHERE mi.inventory_id = inv.id
  	AND m.id = $1;`

	stmt5, err := tx.Prepare(minusInvCycle)
	if err != nil {
		return err
	}
	defer stmt5.Close()

	var wasError, notFound, foundAllergen bool
	var invsTemp []models.InventoryUpdate
	for i, item := range ord.Items {
		var hasInMenu, notEnough bool
		if err = stmt.Get(&ord.Items[i].Allergens, item.ProductID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ord.Items[i].Warning = "not found in menu"
				wasError = true
				notFound = true
			} else {
				return err
			}
		} else if db.checkAllergens(ord.Allergens, &ord.Items[i].Allergens); len(ord.Items[i].Allergens) != 0 {
			ord.Items[i].Warning = "found allergen"
			wasError = true
			foundAllergen = true
		} else if err = stmt2.Select(&ord.Items[i].NotEnoungIngs, item.ProductID, item.Quantity); err != nil {
			return err
		} else if len(ord.Items[i].NotEnoungIngs) != 0 {
			ord.Items[i].Warning = "not enough in inventory"
			wasError = true
			notEnough = true
		} else if !wasError { // insert to order_items
			_, err = stmt3.Exec(ord.ID, item.ProductID, item.Quantity)
			if err != nil {
				return err
			} else if invsUpdatesOriginal != nil {
				// 1 ингридентті 1 тапсырыста 2 меню сол 1еуін қолдануы мүмкін сол үшін керек
				var invsTempTemp []models.InventoryUpdate
				err = stmt4.Select(&invsTempTemp, item.ProductID, item.Quantity)
				if err != nil {
					return err
				}
				db.mergerInv(invsTempTemp, &invsTemp)
			}
		}

		if hasInMenu && !notEnough {
			_, err = stmt5.Exec(item.ProductID, item.Quantity)
			if err != nil {
				return err
			}
		}
	}
	// максимальна клиенттің қатесін басты проритетке аламыз
	if wasError {
		if foundAllergen {
			return models.ErrAllergen // 418 (joke)
		} else if notFound {
			return models.ErrNotFoundItems // 404
		}
		return models.ErrOrderNotEnoughItems // 424
	}

	const queryTransaction string = `
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		SELECT 
			inv.id,
			(ings.quantity * ord.quantity) AS quantity_change,
			'usage'::reason_of_inventory_transaction
		FROM 
			inventory inv
		JOIN 
			menu_item_ingredients ings ON inv.id = ings.inventory_id
		JOIN 
			order_items ord ON ings.product_id = ord.product_id
		WHERE 
			ord.order_id = $1`
	_, err = tx.Exec(queryTransaction, ord.ID)
	if err != nil {
		return err
	}
	// егер бәрі дұрыс болса ғана жолдайды
	if invsUpdatesOriginal != nil {
		db.mergerInv(invsTemp, invsUpdatesOriginal)
	}
	const totalQ string = `
	SELECT SUM(m.price * o.quantity)
	FROM menu_items AS m
	JOIN order_items AS o ON m.id = o.product_id
	WHERE o.order_id = $1`

	ord.Total = new(float64)
	err = tx.Get(ord.Total, totalQ, ord.ID)
	if err != nil {
		return err
	}

	const orderUpdateTotal string = `
	UPDATE orders
	SET total = $1
	WHERE id = $2`
	_, err = tx.Exec(orderUpdateTotal, *ord.Total, ord.ID)
	return err
}

// SELECT inv.id, inv.quantity-(ings.quantity * $1) AS notEnough FROM inventory AS inv JOIN menu_item_ingredients AS ings ON inv.id=ings.inventory_id WHERE ings.product_id = $2 AND inv.quantity-(ings.quantity * $1)<0;

// SELECT inv.id, inv.quantity-(ings.quantity * $1) AS notEnough FROM inventory AS inv JOIN menu_item_ingredients AS ings ON inv.id=ings.inventory_id WHERE ings.product_id = $2 AND inv.quantity-(ings.quantity * $1)<0;
