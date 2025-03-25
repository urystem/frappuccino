package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type MenuRepository struct {
	Db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{
		Db: db,
	}
}

func (r *MenuRepository) GetAll(ctx context.Context) ([]models.Menu, error) {
	query := `
		SELECT *
		FROM menu_items`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []models.Menu
	for rows.Next() {
		var menu models.Menu
		err := rows.Scan(&menu.MenuID, &menu.Name, &menu.Description)
		if err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menu)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return menuItems, nil
}

func (r *MenuRepository) GetElementById(ctx context.Context, menuId int) (models.Menu, error) {
	query := `
		SELECT *
		FROM menu_items
		WHERE menu_item_id = $1`

	var menu models.Menu
	err := r.Db.QueryRowContext(ctx, query, menuId).
		Scan(&menu.Name)

	if err != nil {
		return models.Menu{}, err
	}

	return menu, nil
}

func (r *MenuRepository) Put(ctx context.Context, item models.Menu) error {
	// const op = "repository.menu.Put"

	query := `
		UPDATE menu_items
		SET menu_name = $1, quantity = $2, unit = $3
		WHERE menu_item_id = $4`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, item.Name)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		message := fmt.Sprintf("update failed, menu item with ID %v does not exist", item.MenuID)
		return errors.New(message)
	}

	return nil
}

func (r *MenuRepository) Delete(ctx context.Context, menuId int) error {
	query := `
		DELETE FROM menu_items
		WHERE menu_item_id = $1`

	res, err := r.Db.ExecContext(ctx, query, menuId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		message := fmt.Sprintf("deletion was not successful, menu item with ID %v does not exist", menuId)
		return errors.New(message)
	}

	return nil
}

func (r *MenuRepository) Post(ctx context.Context, item models.Menu) error {
	query := `
		INSERT INTO menu_items (menu_name, quantity, unit) 
		VALUES ($1, $2, $3)`

	_, err := r.Db.ExecContext(ctx, query, item.Name)
	if err != nil {
		return err
	}

	return nil
}
