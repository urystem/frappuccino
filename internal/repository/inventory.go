package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type InventoryRepository struct {
	Db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{Db: db}
}

// GetAll retrieves all inventory items.
func (r *InventoryRepository) GetAll(ctx context.Context) ([]*models.InventoryItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT * FROM inventory_items")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory: %w", err)
	}
	defer rows.Close()

	var inventory []*models.InventoryItem
	for rows.Next() {
		item, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory row: %w", err)
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

// GetByID retrieves an inventory item by ID.
func (r *InventoryRepository) GetByID(ctx context.Context, id int) (*models.InventoryItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT * FROM inventory_items WHERE inventory_item_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory item %d: %w", id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	item, err := scanRowsIntoProduct(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan inventory item %d: %w", id, err)
	}

	return item, nil
}

// Delete removes an inventory item.
func (r *InventoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM inventory_items WHERE inventory_item_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory item %d: %w", id, err)
	}
	return nil
}

// Update modifies an existing inventory item.
func (r *InventoryRepository) Update(ctx context.Context, item *models.InventoryItem) error {
	// Start a transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var prevQuantity int
	// Get the previous quantity
	err = tx.QueryRowContext(ctx, `
    SELECT quantity 
    FROM inventory_items 
    WHERE inventory_item_id = $1`, item.ID).Scan(&prevQuantity)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get previous quantity for item %d: %w", item.ID, err)
	}

	quantityChange := item.Quantity - prevQuantity

	// Update the inventory item
	_, err = tx.ExecContext(ctx, `
    UPDATE inventory_items 
    SET name = $1, quantity = $2, unit = $3, allergens = $4, extra_info = $5 
    WHERE inventory_item_id = $6`,
		item.Name, item.Quantity, item.Unit, item.Allergens, item.ExtraInfo, item.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update inventory item %d: %w", item.ID, err)
	}

	// Insert the quantity change into the transactions table
	_, err = tx.ExecContext(ctx, `
    INSERT INTO inventory_transactions (inventory_item_id, quantity_change) 
    VALUES ($1, $2)`,
		item.ID, quantityChange)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert transaction for item %d: %w", item.ID, err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}

// Insert adds a new inventory item.
func (r *InventoryRepository) Insert(ctx context.Context, item *models.InventoryItem) error {
	_, err := r.Db.ExecContext(ctx, `
		INSERT INTO inventory_items (name, quantity, unit, allergens, extra_info) 
		VALUES ($1, $2, $3, $4, $5)`,
		item.Name, item.Quantity, item.Unit, item.Allergens, item.ExtraInfo)
	if err != nil {
		return fmt.Errorf("failed to insert inventory item %s: %w", item.Name, err)
	}
	return nil
}

// scanRowsIntoProduct scans a database row into an InventoryItem.
func scanRowsIntoProduct(rows *sql.Rows) (*models.InventoryItem, error) {
	item := new(models.InventoryItem)
	if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit, &item.Allergens, &item.ExtraInfo); err != nil {
		return nil, err
	}
	return item, nil
}
