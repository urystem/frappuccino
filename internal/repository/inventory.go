package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type InventoryRepository struct {
	Db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		Db: db,
	}
}

func (r *InventoryRepository) GetAll(ctx context.Context) ([]models.Inventory, error) {
	query := `
		SELECT *
		FROM inventory_items`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventoryItems []models.Inventory
	for rows.Next() {
		var inventory models.Inventory
		err := rows.Scan(&inventory.InventoryId, &inventory.Name, &inventory.Quantity, &inventory.Unit)
		if err != nil {
			return nil, err
		}
		inventoryItems = append(inventoryItems, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inventoryItems, nil
}

func (r *InventoryRepository) GetElementById(ctx context.Context, inventoryId int) (models.Inventory, error) {
	query := `
		SELECT *
		FROM inventory_items
		WHERE inventory_item_id = $1`

	var inventory models.Inventory
	err := r.Db.QueryRowContext(ctx, query, inventoryId).
		Scan(&inventory.InventoryId, &inventory.Name, &inventory.Quantity, &inventory.Unit)

	if err != nil {
		return models.Inventory{}, err
	}

	return inventory, nil
}

func (r *InventoryRepository) Put(ctx context.Context, item models.Inventory) error {
	// const op = "repository.inventory.Put"

	query := `
		UPDATE inventory_items
		SET inventory_name = $1, quantity = $2, unit = $3
		WHERE inventory_item_id = $4`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, item.Name, item.Quantity, item.Unit, item.InventoryId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		message := fmt.Sprintf("update failed, inventory item with ID %v does not exist", item.InventoryId)
		return errors.New(message)
	}

	return nil
}

func (r *InventoryRepository) Delete(ctx context.Context, inventoryId int) error {
	query := `
		DELETE FROM inventory_items
		WHERE inventory_item_id = $1`

	res, err := r.Db.ExecContext(ctx, query, inventoryId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		message := fmt.Sprintf("deletion was not successful, inventory item with ID %v does not exist", inventoryId)
		return errors.New(message)
	}

	return nil
}

func (r *InventoryRepository) Post(ctx context.Context, item models.Inventory) error {
	query := `
		INSERT INTO inventory_items (inventory_name, quantity, unit) 
		VALUES ($1, $2, $3)`

	_, err := r.Db.ExecContext(ctx, query, item.Name, item.Quantity, item.Unit)
	if err != nil {
		return err
	}

	return nil
}
