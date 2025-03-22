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

func (r *InventoryRepository) Put(ctx context.Context, item models.Inventory) error {
	// const op = "repository.inventory.Put"

	query := `
		UPDATE inventory 
		SET inventory_name = $1, quantity = $2, unit = $3, allergens = $4, is_active = $5 
		WHERE inventory_id = $6`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, item.InventoryName, item.Quantity, item.Unit, item.Allergens, item.InventoryId, item.IsActive)
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

func (r *InventoryRepository) Post(ctx context.Context, item models.Inventory) error {
	query := `
		INSERT INTO inventory (inventory_name, quantity, unit, allergens, is_active) 
		VALUES ($1, $2, $3, $4, $5) RETURNING inventory_id`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var lastId int
	err = stmt.QueryRowContext(ctx, item.InventoryName, item.Quantity, item.Unit, item.Allergens, item.IsActive).Scan(&lastId)
	if err != nil {
		return err
	}

	return nil
}

func (r *InventoryRepository) GetAll(ctx context.Context) ([]models.Inventory, error) {
	query := `
		SELECT inventory_id, inventory_name, quantity, unit, allergens, is_active 
		FROM inventory 
		WHERE is_active = TRUE`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventoryItems []models.Inventory
	for rows.Next() {
		var inventory models.Inventory
		err := rows.Scan(&inventory.InventoryId, &inventory.InventoryName, &inventory.Quantity, &inventory.Unit, &inventory.Allergens, &inventory.IsActive)
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
		SELECT inventory_id, inventory_name, quantity, unit, allergens, is_active 
		FROM inventory 
		WHERE inventory_id = $1 AND is_active = TRUE`

	var inventory models.Inventory
	err := r.Db.QueryRowContext(ctx, query, inventoryId).
		Scan(&inventory.InventoryId, &inventory.InventoryName, &inventory.Quantity, &inventory.Unit, &inventory.Allergens, &inventory.IsActive)

	if err != nil {
		return models.Inventory{}, err
	}

	return inventory, nil
}

func (r *InventoryRepository) Delete(ctx context.Context, inventoryId int) error {
	query := `
		UPDATE inventory 
		SET is_active = FALSE 
		WHERE inventory_id = $1 AND is_active = TRUE`

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
