package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

type InventoryDataAccess interface {
	ReadInventory() ([]models.InventoryItem, error)
	WriteInventory([]models.InventoryItem) error
}

type inventoryRepository struct {
	inventFilePath string
}

// Конструктор для InventoryRepository
func NewInventoryRepository(filepath string) *inventoryRepository {
	return &inventoryRepository{inventFilePath: filepath}
}

// Метод для чтения данных инвентаря из файла
func (r *inventoryRepository) ReadInventory() ([]models.InventoryItem, error) {
	file, err := os.Open(r.inventFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var items []models.InventoryItem
	if err = json.NewDecoder(file).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// Метод для записи данных инвентаря в файл
func (r *inventoryRepository) WriteInventory(items []models.InventoryItem) error {
	file, err := os.Create(r.inventFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(items)
}
