package dal

import (
	"database/sql"
	"errors"

	"frappuccino/models"

	"github.com/jmoiron/sqlx"
)

type dalOrder struct {
	database *sqlx.DB
}

type OrderDalInter interface {
	SelectAllOrders() ([]models.Order, error)
	SelectOrder(uint64) (*models.Order, error)
	DeleteOrder(uint64) error
	InsertOrder(*models.Order) error
	UpdateOrder(uint64, *models.Order) error
	CloseOrder(uint64) error
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
	err = db.inventoryUpdaterByOrderAndSetTotal(tx, id, nil)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	// жәй өшіре саламыз. order_items тен өзі өшіп кетеді
	_, err = tx.Exec(`DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

func (db *dalOrder) InsertOrder(ord *models.Order) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = tx.QueryRow(`INSERT INTO orders (customer_name, allergens)
	VALUES($1,$2)
	RETURNING id`, ord.CustomerName, ord.Allergens).Scan(&ord.ID); err != nil {
		return err
	}
	err = db.inventoryUpdaterByOrderAndSetTotal(tx, ord.ID, &ord.Items)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *dalOrder) UpdateOrder(id uint64, ord *models.Order) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = db.inventoryUpdaterByOrderAndSetTotal(tx, id, &ord.Items)
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(`UPDATE orders 
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

func (db *dalOrder) getStatus(tx *sqlx.Tx, id uint64) (string, error) {
	var status string
	err := tx.Get(&status, `SELECT status FROM orders WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return "", models.ErrNotFound
	} else if err != nil {
		return "", err
	}
	return status, nil
}

// delete, insert, update
func (db *dalOrder) inventoryUpdaterByOrderAndSetTotal(tx *sqlx.Tx, id uint64, items *[]models.OrderItem) error {
	// getting a status of order
	status, err := db.getStatus(tx, id)
	if err != nil {
		return err
	}
	// тапсырыс әлі орындалмаған болса, қосамыз
	if status == "processing" {
		// update жасар алдында міндетті түрде селектпен тексеру керек екен)
		_, err := tx.Exec(`UPDATE inventory AS inv
		SET quantity = inv.quantity + (ings.quantity * ord.quantity)
		FROM menu_item_ingredients AS ings
		JOIN order_items AS ord ON ings.product_id = ord.product_id
		WHERE inv.id = ings.inventory_id AND ord.order_id = $1`, id)
		// UPDATE inventory AS inv SET quantity = inv.quantity+(i.quantity * o.quantity) FROM menu_item_ingredients AS i JOIN order_items AS o ON i.product_id = o.product_id WHERE o.order_id = 1;
		if err != nil {
			return err
		}

	} else if items != nil { // егер бұл update болса, әрі ол processing емес болса
		return errors.New("order closed: you cannot update")
	}

	// егер delete ке болса
	if items == nil {
		return nil
	}

	// update үшін тазалап аламыз
	if _, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, id); err != nil {
		return err
	}

	// menuде бар жоғын тексереміз
	// stmt, err := tx.Prepare(`SELECT TRUE FROM menu WHERE id = $1`)
	// EXISTS возвращает true или false
	stmt, err := tx.Prepare(`SELECT EXISTS(SELECT 1 FROM menu_items WHERE id = $1)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	notEnoughInventsQ := `SELECT id, name, ABS(garbage) AS not_enough
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

	stmt3, err := tx.PrepareNamed(`INSERT INTO order_items VALUES(:order_id, :product_id, :quantity)`)
	if err != nil {
		return err
	}
	defer stmt3.Close()

	var wasError bool
	for i, item := range *items {
		var isHasInMenu bool
		if err = stmt.QueryRow(item.ProductID).Scan(&isHasInMenu); err != nil {
			return err
		} else if (*items)[i].NotEnoungIngs = nil; !isHasInMenu {
			(*items)[i].Warning = "not found in menu"
			wasError = true

		} else if err = stmt2.Select(&(*items)[i].NotEnoungIngs, item.ProductID, item.Quantity); err != nil {
			return err
		} else if len((*items)[i].NotEnoungIngs) != 0 {
			(*items)[i].Warning = "not enough in inventory"
			wasError = true
		} else if !wasError { // insert to order_items
			item.OrderId = id
			_, err = stmt3.Exec(item)
			if err != nil {
				return err
			}
		}
	}

	if wasError {
		// *items = slices.DeleteFunc(*items, func(item models.OrderItem) bool {
		// 	return len(item.Err) != 0
		// })
		var invalidItemsCount uint64
		for _, item := range *items {
			if len(item.Warning) != 0 {
				(*items)[invalidItemsCount] = item
				invalidItemsCount++
			}
		}
		*items = (*items)[:invalidItemsCount]
		return models.ErrOrderItems
	}

	// updater inventory
	_, err = tx.Exec(`UPDATE inventory AS inv
		SET quantity = inv.quantity - (ings.quantity * ord.quantity)
		FROM menu_item_ingredients AS ings
		JOIN order_items AS ord ON ings.product_id = ord.product_id
		WHERE inv.id = ings.inventory_id AND ord.order_id = $1`, id)
	if err != nil {
		return err
	}

	totalQ := `
	SELECT SUM(m.price * o.quantity)
	FROM menu_items AS m
	JOIN order_items AS o ON m.id = o.product_id
	WHERE o.order_id = $1`

	var total float64
	err = tx.Get(&total, totalQ, id)
	if err != nil {
		return err
	}

	orderUpdateTotal := `UPDATE orders
	SET total = $1
	WHERE id = $2`

	res, err := tx.Exec(orderUpdateTotal, total, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("no rows were updated — order with id not found")
	}
	return nil
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
		return models.ErrOrdStatusClosed
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

// SELECT inv.id, inv.quantity-(ings.quantity * $1) AS notEnough FROM inventory AS inv JOIN menu_item_ingredients AS ings ON inv.id=ings.inventory_id WHERE ings.product_id = $2 AND inv.quantity-(ings.quantity * $1)<0;

// SELECT inv.id, inv.quantity-(ings.quantity * $1) AS notEnough FROM inventory AS inv JOIN menu_item_ingredients AS ings ON inv.id=ings.inventory_id WHERE ings.product_id = $2 AND inv.quantity-(ings.quantity * $1)<0;
