package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
)

type InventoryRepository struct {
	Db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		Db: db,
	}
}

// Get all inventory items
func (r *InventoryRepository) GetAll(ctx context.Context) ([]*models.InventoryItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT * FROM inventory_items")
	if err != nil {
		return nil, err
	}

	inventory := make([]*models.InventoryItem, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		inventory = append(inventory, p)
	}

	return inventory, nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, id int) (*models.InventoryItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT * FROM inventory_items WHERE inventory_item_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	item, err := scanRowsIntoProduct(rows)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// Delete an inventory item
func (r *InventoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM inventory_items WHERE inventory_item_id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// Update an inventory item
func (r *InventoryRepository) Update(ctx context.Context, item *models.InventoryItem) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE inventory_items SET name=$1, quantity=$2, unit=$3, allergens=$4, extra_info=$5 where inventory_item_id=$6", item.Name, item.Quantity, item.Unit, item.Allergens, item.ExtraInfo, item.ID)
	if err != nil {
		return err
	}
	return nil
}

// Insert a new inventory item
func (r *InventoryRepository) Insert(ctx context.Context, item *models.InventoryItem) error {
	_, err := r.Db.ExecContext(ctx, "INSERT INTO inventory_items (name, quantity, unit, allergens, extra_info) VALUES ($1, $2, $3, $4, $5)", item.Name, item.Quantity, item.Unit, item.Allergens, item.ExtraInfo)
	if err != nil {
		return err
	}
	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*models.InventoryItem, error) {
	inventory := new(models.InventoryItem)

	err := rows.Scan(
		&inventory.ID,
		&inventory.Name,
		&inventory.Quantity,
		&inventory.Unit,
		&inventory.Allergens,
		&inventory.ExtraInfo,
	)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}
