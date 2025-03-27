package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
)

type MenuRepository struct {
	Db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{Db: db}
}

func (r *MenuRepository) GetAll(ctx context.Context) ([]*models.MenuItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT menu_item_id, name, description, details, price, allergens FROM menu_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menuItems := make([]*models.MenuItem, 0)
	for rows.Next() {
		item := new(models.MenuItem)
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Details, &item.Price, &item.Allergens)
		if err != nil {
			return nil, err
		}

		item.Ingredients, err = r.getIngredients(ctx, item.ID)
		if err != nil {
			return nil, err
		}

		menuItems = append(menuItems, item)
	}

	return menuItems, nil
}

func (r *MenuRepository) GetByID(ctx context.Context, id int) (*models.MenuItem, error) {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "SELECT menu_item_id, name, description, details, price, allergens FROM menu_items WHERE menu_item_id = $1", id)

	item := new(models.MenuItem)
	err = row.Scan(&item.ID, &item.Name, &item.Description, &item.Details, &item.Price, &item.Allergens)
	if err != nil {
		return nil, err
	}

	item.Ingredients, err = r.getIngredients(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *MenuRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM menu_item_ingredients WHERE menu_item_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM menu_items WHERE menu_item_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *MenuRepository) Update(ctx context.Context, item *models.MenuItem) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var oldPrice float64
	err = tx.QueryRowContext(ctx, "SELECT price FROM menu_items WHERE menu_item_id = $1", item.ID).Scan(&oldPrice)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE menu_items SET name = $1, description = $2, price = $3, allergens = $4 WHERE menu_item_id = $5",
		item.Name, item.Description, item.Price, item.Allergens, item.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if oldPrice != item.Price {
		_, err = tx.ExecContext(ctx, "INSERT INTO price_history (menu_item_id, old_price, new_price) VALUES ($1, $2, $3)",
			item.ID, oldPrice, item.Price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = r.updateIngredients(ctx, item.ID, item.Ingredients)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *MenuRepository) Insert(ctx context.Context, item *models.MenuItem) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var menuItemID int
	err = tx.QueryRowContext(ctx, "INSERT INTO menu_items (name, description, details, price, allergens) VALUES ($1, $2, $3, $4, $5) RETURNING menu_item_id",
		item.Name, item.Description, item.Details, item.Price, item.Allergens).Scan(&menuItemID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, ingredient := range item.Ingredients {
		_, err := tx.ExecContext(ctx, "INSERT INTO menu_item_ingredients (menu_item_id, inventory_item_id, quantity) VALUES ($1, $2, $3)",
			menuItemID, ingredient.InventoryItemID, ingredient.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *MenuRepository) getIngredients(ctx context.Context, menuItemID int) ([]models.MenuItemIngredient, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT menu_item_ingredient_id, menu_item_id, inventory_item_id, quantity FROM menu_item_ingredients WHERE menu_item_id = $1", menuItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := make([]models.MenuItemIngredient, 0)
	for rows.Next() {
		ingredient := models.MenuItemIngredient{}
		err := rows.Scan(&ingredient.ID, &ingredient.MenuItemID, &ingredient.InventoryItemID, &ingredient.Quantity)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (r *MenuRepository) updateIngredients(ctx context.Context, menuItemID int, ingredients []models.MenuItemIngredient) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM menu_item_ingredients WHERE menu_item_id = $1", menuItemID)
	if err != nil {
		return err
	}

	for _, ingredient := range ingredients {
		ingredient.MenuItemID = menuItemID
		if err := r.insertIngredient(ctx, ingredient); err != nil {
			return err
		}
	}

	return nil
}

func (r *MenuRepository) insertIngredient(ctx context.Context, ingredient models.MenuItemIngredient) error {
	_, err := r.Db.ExecContext(ctx, "INSERT INTO menu_item_ingredients (menu_item_id, inventory_item_id, quantity) VALUES ($1, $2, $3) RETURNING menu_item_ingredient_id",
		ingredient.MenuItemID, ingredient.InventoryItemID, ingredient.Quantity)
	return err
}
