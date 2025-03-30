package models

import (
	"errors"
	"time"
)

type OrderStatus int

const (
	OrderPending OrderStatus = iota
	OrderConfirmed
	OrderInProgress
	OrderCompleted
	OrderCancelled
	OrderRejected
)

func (s OrderStatus) String() string {
	switch s {
	case OrderPending:
		return "pending"
	case OrderConfirmed:
		return "confirmed"
	case OrderInProgress:
		return "in progress"
	case OrderCompleted:
		return "completed"
	case OrderCancelled:
		return "cancelled"
	case OrderRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderPending, OrderConfirmed, OrderInProgress, OrderCompleted, OrderCancelled, OrderRejected:
		return true
	}
	return false
}

type Order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customer_name"`
	Status       string      `json:"status"`
	Total        float32     `json:"total"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Items        []OrderItem `json:"items"`
}

type OrderItem struct {
	ID         int `json:"id"`
	OrderID    int `json:"order_id"`
	MenuItemID int `json:"menu_item_id"`
	Quantity   int `json:"quantity"`
}

type OrderStatusHistory struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (o *OrderItem) IsValid() error {
	if o.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return nil
}
