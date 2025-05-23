package dal

import (
	"database/sql"

	"frappuccino/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type dalInv struct {
	db *sqlx.DB
}

type InventoryDataAccess interface {
	// InsertInventoryV1(*models.Inventory) error
	// InsertInventoryV2(*models.Inventory) error
	// InsertInventoryV3(*models.Inventory) error
	// InsertInventoryV4(*models.Inventory) error
	InsertInventoryV5(*models.Inventory) error
	// InsertInventoryV6(*models.Inventory) error
	SelectAllInventories() ([]models.Inventory, error)
	SelectInventory(uint64) (*models.Inventory, error)
	UpdateInventory(*models.Inventory) error
	DeleteInventory(uint64) (*models.InventoryDepend, error)
	SelectAllInventoryTransaction() ([]models.InventoryTransaction, error)
	SelectReorder() ([]models.Inventory, error)
}

func ReturnDalInvCore(db *sqlx.DB) InventoryDataAccess {
	return &dalInv{db: db}
}

func (core *dalInv) InsertInventoryV1(inv *models.Inventory) error {
	_, err := core.db.Exec(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES($1,$2,$3,$4,$5,$6)`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price,
	)
	return err
}

func (core *dalInv) InsertInventoryV2(inv *models.Inventory) error {
	_, err := core.db.NamedExec(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES (:name, :description, :quantity, :reorder_level, :unit, :price)
	`, inv)
	return err
}

func (core *dalInv) InsertInventoryV3(inv *models.Inventory) error {
	return core.db.QueryRow(`
	INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
		VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price).Scan(&inv.ID)
}

func (core *dalInv) InsertInventoryV4(inv *models.Inventory) error {
	return core.db.Get(&inv.ID, `
    INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id`, inv.Name, inv.Descrip, inv.Quantity, inv.ReorderLvl, inv.Unit, inv.Price)
}

func (core *dalInv) InsertInventoryV5(inv *models.Inventory) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	// var ss any

	defer tx.Rollback()
	// tx.QueryRowx также подходит
	if err = tx.QueryRow(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price).Scan(&inv.ID); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique
				return models.ErrConflict
			}
		}
		return err
	}
	_, err = tx.NamedExec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES (:id, :quantity, 'restock')
	`, inv)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalInv) InsertInventoryV6(inv *models.Inventory) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareNamed(`
	INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
	VALUES (:name, :description, :quantity, :reorder_level, :unit, :price)
	RETURNING id`)
	if err != nil {
		return err // Ошибка при подготовке запроса
	}
	defer stmt.Close()

	// Выполнение подготовленного запроса
	if err = stmt.QueryRow(inv).Scan(&inv.ID); err != nil {
		return err // Ошибка при выполнении
	}

	if _, err = tx.NamedExec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES (:id, :quantity, 'restock')
	`, inv); err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalInv) SelectAllInventories() ([]models.Inventory, error) {
	var invts []models.Inventory
	return invts, core.db.Select(&invts, "SELECT * FROM inventory")
}

func (core *dalInv) SelectInventory(id uint64) (*models.Inventory, error) {
	var inv models.Inventory
	err := core.db.Get(&inv, "SELECT * FROM inventory WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	}
	return &inv, err
}

func (core *dalInv) UpdateInventory(inv *models.Inventory) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var quantity_changed float64
	// err = tx.QueryRow(`SELECT quantity FROM inventory WHERE id=$1`,inv.ID).Scan(&oldQuantity)
	if err = tx.Get(&quantity_changed, `SELECT quantity FROM inventory WHERE id=$1`, inv.ID); err != nil {
		if err == sql.ErrNoRows {
			return models.ErrNotFound
		}
		return err
	}

	_, err = tx.NamedExec(`
	UPDATE inventory
		SET name = :name, description = :description, quantity = :quantity,
		    reorder_level = :reorder_level, unit = :unit, price = :price
		WHERE id = :id`, inv)
	if err != nil {
		return err
	}

	quantity_changed = inv.Quantity - quantity_changed

	reason := "restock"
	if quantity_changed < 0 {
		reason = "annul"
	}

	if _, err = tx.Exec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES ($1, $2, $3)
	`, inv.ID, quantity_changed, reason); err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalInv) DeleteInventory(id uint64) (*models.InventoryDepend, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var menuDepend models.InventoryDepend

	const menusNames string = `SELECT id, name
		FROM menu_items
		JOIN menu_item_ingredients ON id=product_id
		WHERE inventory_id=$1`

	err = tx.Select(&menuDepend.Menus, menusNames, id)
	if err != nil {
		return nil, err
	}

	if len(menuDepend.Menus) != 0 {
		return &menuDepend, nil
	}
	res, err := tx.Exec(`DELETE FROM inventory WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	affects, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affects == 0 {
		return nil, models.ErrNotFound
	}
	return nil, tx.Commit()
}

func (core *dalInv) SelectAllInventoryTransaction() ([]models.InventoryTransaction, error) {
	var inventoryTransactions []models.InventoryTransaction
	err := core.db.Select(&inventoryTransactions, `SELECT * FROM inventory_transactions ORDER BY updated_at ASC`)
	if err != nil {
		return nil, err
	}
	return inventoryTransactions, nil
}

func (core *dalInv) SelectReorder() ([]models.Inventory, error) {
	var invs []models.Inventory
	return invs, core.db.Select(&invs, `SELECT * FROM inventory WHERE quantity <= reorder_level`)
}
