package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

// GetAll retrieves all orders along with their items and status history.
func (r *OrderRepository) GetAll(ctx context.Context) ([]*models.Order, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT order_id, customer_name, status, total FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := new(models.Order)
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.Status, &order.Total); err != nil {
			return nil, err
		}

		// Fetch order items
		order.Items, err = r.getOrderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		// Fetch status history and assign the latest status
		history, err := r.getOrderStatusHistory(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		if len(history) > 0 {
			order.Status = history[len(history)-1].Status.String()
		}

		orders = append(orders, order)
	}
	return orders, nil
}

// GetByID retrieves a single order by its ID along with its items and status history.
func (r *OrderRepository) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	order := new(models.Order)

	// Fetch order details
	err := r.Db.QueryRowContext(ctx, "SELECT order_id, customer_name, status, total FROM orders WHERE order_id = $1",
		orderID).Scan(&order.ID, &order.CustomerName, &order.Status, &order.Total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Fetch order items
	order.Items, err = r.getOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	// Fetch status history
	history, err := r.getOrderStatusHistory(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	// Assign the latest status if history exists
	if len(history) > 0 {
		order.Status = history[len(history)-1].Status.String()
	}

	return order, nil
}

// Insert creates a new order while checking allergens and inventory.
func (r *OrderRepository) Insert(ctx context.Context, order *models.Order) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check allergens
	for _, item := range order.Items {
		if err := r.checkForAllergens(ctx, tx, order.ID, item.ID); err != nil {
			return err
		}
	}

	// Check inventory
	for _, item := range order.Items {
		if err := r.checkInventory(ctx, tx, item.ID, item.Quantity); err != nil {
			return err
		}
	}

	// Insert order
	var orderID int
	err = tx.QueryRowContext(ctx, "INSERT INTO orders (customer_name, status, total) VALUES ($1, 'Pending', $2) RETURNING order_id",
		order.CustomerName, order.Total).Scan(&orderID)
	if err != nil {
		return err
	}

	// Insert order items
	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, "INSERT INTO order_items (order_id, item_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.ID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Insert initial status history
	_, err = tx.ExecContext(ctx, "INSERT INTO order_status_history (order_id, status, updated_at) VALUES ($1, 'Pending', $2)",
		orderID, time.Now())
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// UpdateStatus updates the order status and records history.
func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE order_id = $2", order.Status, order.ID)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, "INSERT INTO order_status_history (order_id, status, updated_at) VALUES ($1, $2, $3)",
		order.ID, order.Status, time.Now())
	return err
}

// Delete removes an order and its related records (items and history).
func (r *OrderRepository) Delete(ctx context.Context, orderID int) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete order items
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}

	// Delete order status history
	_, err = tx.ExecContext(ctx, "DELETE FROM order_status_history WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}

	// Delete order
	res, err := tx.ExecContext(ctx, "DELETE FROM orders WHERE order_id = $1", orderID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return tx.Commit()
}

// getOrderItems retrieves the items for a given order.
func (r *OrderRepository) getOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT item_id, quantity FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// getOrderStatusHistory retrieves status change history.
func (r *OrderRepository) getOrderStatusHistory(ctx context.Context, orderID int) ([]models.OrderStatusHistory, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT status, updated_at FROM order_status_history WHERE order_id = $1 ORDER BY updated_at ASC", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderStatusHistory
	for rows.Next() {
		var record models.OrderStatusHistory
		if err := rows.Scan(&record.Status, &record.UpdatedAt); err != nil {
			return nil, err
		}
		history = append(history, record)
	}
	return history, nil
}

// checkForAllergens ensures an order does not contain allergens for a customer.
func (r *OrderRepository) checkForAllergens(ctx context.Context, tx *sql.Tx, customerID, itemID int) error {
	var count int
	err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM allergens WHERE customer_id = $1 AND ingredient_id IN (SELECT ingredient_id FROM item_ingredients WHERE item_id = $2)", customerID, itemID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("order contains allergens")
	}
	return nil
}

// checkInventory ensures enough ingredients are available.
func (r *OrderRepository) checkInventory(ctx context.Context, tx *sql.Tx, itemID, quantity int) error {
	var available int
	err := tx.QueryRowContext(ctx, "SELECT stock FROM inventory WHERE item_id = $1", itemID).Scan(&available)
	if err != nil {
		return err
	}
	if available < quantity {
		return errors.New("not enough stock")
	}
	return nil
}
