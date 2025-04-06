package dal

import (
	"database/sql"
	"fmt"

	"frappuccino/models"

	"github.com/jmoiron/sqlx"
)

type OrderDalInter interface {
	SelectAllOrders() ([]models.Order, error)
	SelectOrder(uint64) (*models.Order, error)
	InsertOrder(*models.Order) ([]models.OrderItem, error)
	DeleteOrder(uint64) error
}

func (core *dalCore) SelectAllOrders() ([]models.Order, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orders []models.Order

	err = tx.Select(&orders, `SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareNamed(`SELECT product_id, quantity FROM order_items WHERE order_id = :id`)
	if err != nil {
		return nil, nil
	}

	defer stmt.Close()

	for i, v := range orders {
		err = stmt.Select(&orders[i].Items, v)
		if err != nil {
			return nil, err
		}
	}
	return orders, tx.Commit()
}

func (core *dalCore) SelectOrder(id uint64) (*models.Order, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var order models.Order

	err = tx.Get(&order, `SELECT * FROM orders WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	err = tx.Select(&order.Items, `SELECT product_id, quantity FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &order, tx.Commit()
}

func (core *dalCore) DeleteOrder(id uint64) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.Get(&status, `SELECT status FROM orders WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return models.ErrNotFound
	} else if err != nil {
		return err
	}

	if status == "processing" {
		_, err = tx.Exec(`UPDATE inventory 
		SET quantity = quantity+(i.quantity*o.quantity) 
		FROM menu_item_ingredients AS i
		JOIN order_items AS o ON i.product_id = o.product_id
		WHERE o.order_id = $1`, id)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(`DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalCore) inventoryUpdaterByOrder(tx *sqlx.Tx, items []models.OrderItem) error {
	stmt, err := tx.PrepareNamed(`SELECT inventory_id, quantity*:quantity FROM menu_item_ingredients WHERE product_id = :product_id`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var menuIngs []models.MenuIngredients
	for _, v := range items {
		var menuIngsTemp []models.MenuIngredients
		stmt.Select(&menuIngsTemp, v)
		menuIngs = append(menuIngs, menuIngsTemp...)
	}
	fmt.Println(menuIngs)
	return nil
}

func (core *dalCore) InsertOrder(*models.Order) ([]models.OrderItem, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return nil, tx.Commit()
}

// `UPDATE inventory
// 	SET quantity = quantity+(i.quantity*o.quantity)
// 	FROM menu_item_ingredients AS i
// 	JOIN order_items AS o ON i.product_id = o.product_id
// 	WHERE o.order_id = $1`
