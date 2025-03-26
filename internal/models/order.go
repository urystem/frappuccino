package models

import "time"

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

type OrderItem struct {
	ID         int `json:"id"`
	OrderID    int `json:"order_id"`
	MenuItemID int `json:"menu_item_id"`
	Quantity   int `json:"quantity"`
}

type OrderStatusHistory struct {
	ID        int         `json:"id"`
	OrderID   int         `json:"order_id"`
	Status    OrderStatus `json:"status"`
	UpdatedAt time.Time   `json:"updated_at"`
}
