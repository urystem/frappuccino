package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]*models.Order, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT order_id, customer_name, status, total FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*models.Order, 0)
	for rows.Next() {
		order := new(models.Order)
		err := rows.Scan(&order.ID, &order.CustomerName, &order.Status, &order.Total)
		if err != nil {
			return nil, err
		}

		order.Items, err = r.getOrderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*models.Order, error) {
	row := r.Db.QueryRowContext(ctx, "SELECT order_id, customer_name, status, total FROM orders WHERE order_id = $1", id)

	order := new(models.Order)
	err := row.Scan(&order.ID, &order.CustomerName, &order.Status, &order.Total)
	if err != nil {
		return nil, err
	}

	order.Items, err = r.getOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM orders WHERE order_id = $1", id)
	return err
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE orders SET customer_name = $1, status = $2, total = $3 WHERE order_id = $4",
		order.CustomerName, order.Status, order.Total, order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.updateOrderItems(ctx, tx, order.ID, order.Items)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) Insert(ctx context.Context, order *models.Order) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var orderID int
	err = tx.QueryRowContext(ctx, "INSERT INTO orders (customer_name, status, total) VALUES ($1, $2, $3) RETURNING order_id",
		order.CustomerName, order.Status, order.Total).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, "INSERT INTO order_items (order_id, menu_item_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.MenuItemID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) getOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT order_item_id, order_id, menu_item_id, quantity FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.OrderItem, 0)
	for rows.Next() {
		item := models.OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.MenuItemID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *OrderRepository) updateOrderItems(ctx context.Context, tx *sql.Tx, orderID int, items []models.OrderItem) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err := tx.ExecContext(ctx, "INSERT INTO order_items (order_id, menu_item_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.MenuItemID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}
