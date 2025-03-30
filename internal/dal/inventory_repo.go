package dal

import (
	"database/sql"

	"hot-coffee/models"
)

type InventoryDataAccess interface {
	InsertInventory(*models.InventoryItem) (uint, error)
	ReadInventory() ([]models.InventoryItem, error)
}

type inventoryRepository struct {
	db *sql.DB
}

// Конструктор для InventoryRepository
func NewInventoryRepository(arg_db *sql.DB) *inventoryRepository {
	return &inventoryRepository{db: arg_db}
}

func (invCore *inventoryRepository) InsertInventory(inv *models.InventoryItem) (uint, error) {
	tx, err := invCore.db.Begin()
	if err != nil {
		return 0, err
	}
	var lastId uint
	_, err = tx.Exec(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES($1,$2,$3,$4,$5,$6)`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price,
	)
	return lastId, txAfter(tx, err)
}

// Метод для чтения данных инвентаря из файла
func (r *inventoryRepository) ReadInventory() ([]models.InventoryItem, error) {
	// file, err := os.Open(r.inventFilePath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()
	var items []models.InventoryItem
	// if err = json.NewDecoder(file).Decode(&items); err != nil {
	// 	return nil, err
	// }
	return items, nil
}

// Метод для записи данных инвентаря в файл
func (r *inventoryRepository) WriteInventory(items []models.InventoryItem) error {
	// file, err := os.Create(r.inventFilePath)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	// encoder := json.NewEncoder(file)
	// encoder.SetIndent("", " ")
	return nil
}

func txAfter(tx *sql.Tx, err error) error {
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
