package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	rows, err := r.Db.QueryContext(ctx, "SELECT orders_id, customer_name, status, total FROM orders")
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
			order.Status = history[len(history)-1].Status
		}

		orders = append(orders, order)
	}
	return orders, nil
}

// GetByID retrieves a single order by its ID along with its items and status history.
func (r *OrderRepository) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	order := new(models.Order)

	// Fetch order details
	err := r.Db.QueryRowContext(ctx, "SELECT orders_id, customer_name, status, total FROM orders WHERE orders_id = $1",
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
		order.Status = history[len(history)-1].Status
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

	// Check inventory
	for _, item := range order.Items {
		if err := r.checkInventory(ctx, tx, item.ID, item.Quantity); err != nil {
			return err
		}
	}

	// Insert order
	var orderID int
	err = tx.QueryRowContext(ctx, "INSERT INTO orders (customer_name, status, total) VALUES ($1, 'Pending', $2) RETURNING orders_id",
		order.CustomerName, order.Total).Scan(&orderID)
	if err != nil {
		return err
	}

	// Insert order items
	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, "INSERT INTO order_items (orders_id, menu_items_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.ID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Insert initial status history
	_, err = tx.ExecContext(ctx, "INSERT INTO order_status_history (orders_id, status, updated_at) VALUES ($1, 'Pending', $2)",
		orderID, time.Now())
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// UpdateStatus updates the order status and records history.
func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE orders_id = $2", order.Status, order.ID)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, "INSERT INTO order_status_history (orders_id, status, updated_at) VALUES ($1, $2, $3)",
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
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE orders_id = $1", orderID)
	if err != nil {
		return err
	}

	// Delete order status history
	_, err = tx.ExecContext(ctx, "DELETE FROM order_status_history WHERE orders_id = $1", orderID)
	if err != nil {
		return err
	}

	// Delete order
	res, err := tx.ExecContext(ctx, "DELETE FROM orders WHERE orders_id = $1", orderID)
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

// repository/order_repository.go
func (r *OrderRepository) ProcessBatchOrders(ctx context.Context, orders []models.BatchOrder) ([]models.OrderResult, models.BatchSummary, error) {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, models.BatchSummary{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var results []models.OrderResult
	var summary models.BatchSummary
	summary.TotalOrders = len(orders)
	inventoryUpdates := make(map[int]models.InventoryUpdate)

	for _, order := range orders {
		result := models.OrderResult{
			CustomerName: order.CustomerName,
		}

		// Check inventory first
		canFulfill, inventoryCheck := r.checkInventoryBatch(ctx, tx, order.Items)
		if !canFulfill {
			result.Status = "rejected"
			result.Reason = "insufficient_inventory"
			results = append(results, result)
			summary.Rejected++
			continue
		}

		// Create order
		orderID, total, err := r.createOrder(ctx, tx, order)
		if err != nil {
			return nil, models.BatchSummary{}, fmt.Errorf("failed to create order: %w", err)
		}

		// Update inventory
		for _, update := range inventoryCheck {
			if existing, ok := inventoryUpdates[update.IngredientID]; ok {
				existing.QuantityUsed += update.QuantityUsed
				existing.Remaining = update.Remaining
				inventoryUpdates[update.IngredientID] = existing
			} else {
				inventoryUpdates[update.IngredientID] = update
			}
		}

		result.OrderID = orderID
		result.Status = "accepted"
		result.Total = total
		results = append(results, result)
		summary.Accepted++
		summary.TotalRevenue += total
	}

	// Convert inventory updates to slice
	for _, update := range inventoryUpdates {
		summary.InventoryUpdates = append(summary.InventoryUpdates, update)
	}

	if err := tx.Commit(); err != nil {
		return nil, models.BatchSummary{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return results, summary, nil
}

// Update the checkInventory method in repository/order_repository.go
func (r *OrderRepository) checkInventoryBatch(ctx context.Context, tx *sql.Tx, items []models.OrderItem) (bool, []models.InventoryUpdate) {
	var updates []models.InventoryUpdate

	for _, item := range items {
		// Get ingredients for menu item with names
		var ingredients []struct {
			ID       int
			Name     string
			Quantity int
		}
		rows, err := tx.QueryContext(ctx, `
            SELECT mi.inventory_items_id, ii.name, mi.quantity 
            FROM menu_item_ingredients mi
            JOIN inventory_items ii ON mi.inventory_items_id = ii.inventory_items_id
            WHERE mi.menu_items_id = $1`, item.MenuItemID)
		if err != nil {
			return false, nil
		}
		defer rows.Close()

		for rows.Next() {
			var ing struct {
				ID       int
				Name     string
				Quantity int
			}
			if err := rows.Scan(&ing.ID, &ing.Name, &ing.Quantity); err != nil {
				return false, nil
			}
			ingredients = append(ingredients, ing)
		}

		// Check each ingredient
		for _, ing := range ingredients {
			var currentStock int
			err := tx.QueryRowContext(ctx, `
                SELECT quantity 
                FROM inventory_items 
                WHERE inventory_items_id = $1 FOR UPDATE`, ing.ID).Scan(&currentStock)
			if err != nil {
				return false, nil
			}

			needed := ing.Quantity * item.Quantity
			if currentStock < needed {
				return false, nil
			}

			updates = append(updates, models.InventoryUpdate{
				IngredientID: ing.ID,
				Name:         ing.Name,
				QuantityUsed: needed,
				Remaining:    currentStock - needed,
			})
		}
	}

	return true, updates
}

func (r *OrderRepository) createOrder(ctx context.Context, tx *sql.Tx, order models.BatchOrder) (int, float64, error) {
	var orderID int
	var total float64

	// Calculate total
	for _, item := range order.Items {
		var price float64
		err := tx.QueryRowContext(ctx, `
            SELECT price FROM menu_items WHERE menu_items_id = $1`, item.MenuItemID).Scan(&price)
		if err != nil {
			return 0, 0, err
		}
		total += price * float64(item.Quantity)
	}

	// Create order
	err := tx.QueryRowContext(ctx, `
        INSERT INTO orders (customer_name, total, status)
        VALUES ($1, $2, 'pending')
        RETURNING orders_id`, order.CustomerName, total).Scan(&orderID)
	if err != nil {
		return 0, 0, err
	}

	// Add order items
	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, `
            INSERT INTO order_items (orders_id, menu_items_id, quantity)
            VALUES ($1, $2, $3)`, orderID, item.MenuItemID, item.Quantity)
		if err != nil {
			return 0, 0, err
		}
	}

	// Update inventory
	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, `
            UPDATE inventory_items ii
            SET quantity = ii.quantity - (mii.quantity * $1)
            FROM menu_item_ingredients mii
            WHERE mii.menu_items_id = $2
            AND ii.inventory_items_id = mii.inventory_items_id`, item.Quantity, item.MenuItemID)
		if err != nil {
			return 0, 0, err
		}
	}

	// Update order status
	_, err = tx.ExecContext(ctx, `
        INSERT INTO order_status_history (orders_id, status)
        VALUES ($1, 'completed')`, orderID)
	if err != nil {
		return 0, 0, err
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE orders SET status = 'completed' WHERE orders_id = $1`, orderID)
	if err != nil {
		return 0, 0, err
	}

	return orderID, total, nil
}

// getOrderItems retrieves the items for a given order.
func (r *OrderRepository) getOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT menu_items_id, quantity FROM order_items WHERE orders_id = $1", orderID)
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
	rows, err := r.Db.QueryContext(ctx, "SELECT status, updated_at FROM order_status_history WHERE orders_id = $1 ORDER BY updated_at ASC", orderID)
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

// checkInventory ensures enough ingredients are available.
func (r *OrderRepository) checkInventory(ctx context.Context, tx *sql.Tx, itemID, quantity int) error {
	var available int
	err := tx.QueryRowContext(ctx, "SELECT stock FROM inventory WHERE menu_items_id = $1", itemID).Scan(&available)
	if err != nil {
		return err
	}
	if available < quantity {
		return errors.New("not enough stock")
	}
	return nil
}
