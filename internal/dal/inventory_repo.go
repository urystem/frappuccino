package dal

import (
	"hot-coffee/models"

	"github.com/jmoiron/sqlx"
)

type InventoryDataAccess interface {
	// InsertInventoryV1(*models.Inventory) error
	// InsertInventoryV2(*models.Inventory) error
	// InsertInventoryV3(*models.Inventory) error
	// InsertInventoryV4(*models.Inventory) error
	InsertInventoryV5(*models.Inventory) error
	// InsertInventoryV6(*models.Inventory) error
	GetAllInventory() ([]models.Inventory, error)
	SelectInventory(uint64) (*models.Inventory, error)
	UpdateInventory(inv *models.Inventory) error
}

type inventoryRepository struct {
	db *sqlx.DB
}

// Конструктор для InventoryRepository
func NewInventoryRepository(arg_db *sqlx.DB) *inventoryRepository {
	return &inventoryRepository{db: arg_db}
}

func (invCore *inventoryRepository) InsertInventoryV1(inv *models.Inventory) error {
	_, err := invCore.db.Exec(`
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

func (invCore *inventoryRepository) InsertInventoryV2(inv *models.Inventory) error {
	_, err := invCore.db.NamedExec(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES (:name, :description, :quantity, :reorder_level, :unit, :price)
	`, inv)
	return err
}

func (invCore *inventoryRepository) InsertInventoryV3(inv *models.Inventory) error {
	return invCore.db.QueryRow(`
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

func (invCore *inventoryRepository) InsertInventoryV4(inv *models.Inventory) error {
	return invCore.db.Get(&inv.ID, `
    INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id`, inv.Name, inv.Descrip, inv.Quantity, inv.ReorderLvl, inv.Unit, inv.Price)
}

func (invCore *inventoryRepository) InsertInventoryV5(inv *models.Inventory) error {
	tx, err := invCore.db.Beginx()
	if err != nil {
		return err
	}
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
		return err
	} else if _, err = tx.NamedExec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES (:id, :quantity, 'restock')
	`, inv); err != nil {
		return err
	}
	return tx.Commit()
}

func (invCore *inventoryRepository) InsertInventoryV6(inv *models.Inventory) error {
	tx, err := invCore.db.Beginx()
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
	} else if _, err = tx.NamedExec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES (:id, :quantity, 'restock')
	`, inv); err != nil {
		return err
	}
	return tx.Commit()
}

func (invCore *inventoryRepository) GetAllInventory() ([]models.Inventory, error) {
	var invts []models.Inventory
	err := invCore.db.Select(&invts, "SELECT * FROM inventory")
	return invts, err
}

func (invCore *inventoryRepository) SelectInventory(id uint64) (*models.Inventory, error) {
	var inv models.Inventory
	err := invCore.db.Get(&inv, "SELECT * FROM inventory WHERE id = $1", id)
	return &inv, err
}

func (invCore *inventoryRepository) UpdateInventory(inv *models.Inventory) error {
	tx, err := invCore.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var quantity_changed float64
	// err = tx.QueryRow(`SELECT quantity FROM inventory WHERE id=$1`,inv.ID).Scan(&oldQuantity)
	if err = tx.Get(&quantity_changed, `SELECT quantity FROM inventory WHERE id=$1`, inv.ID); err != nil {
		return err
	} else if res, err := tx.NamedExec(`
	UPDATE inventory
		SET name = :name, description = :description, quantity = :quantity,
		    reorder_level = :reorder_level, unit = :unit, price = :price
		WHERE id = :id
	`, inv); err != nil {
		return err
		// UPDATE егер id табылмаса да қате қайтармайды.
		// сондықтан кестенің өзгерісін rowsAffected пен тексереміз
		// егер айди бар болып, бірақ жаңа кестеден айырмасы жоқ болса да, rowsAffected =1 болады
		// демек тек жоқ айдиде ғана 0 болады
	} else if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if quantity_changed = inv.Quantity - quantity_changed; rowsAffected == 0 {
		return models.ErrNotFound
	}
	reason := "restock"
	if quantity_changed < 0 {
		reason = "annul"
	}
	if _, err = tx.Exec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES ($1, $2, $3)
	`, inv.ID, inv.Quantity); err != nil {
		return err
	}
	return tx.Commit()
}
