package models

import (
	"time"
)

type Order struct {
	OrderID     int         `json:"order_id"`
	CustomerID  int         `json:"customer_id"`
	CreatedAt   time.Time   `json:"created_at"`
	Status      OrderStatus `json:"status"`
	TotalAmount int         `json:"total_amount"`
}

type OrderStatusHistory struct {
	OrderHistoryID int         `json:"order_history_id"`
	OrderID        int         `json:"order_id"`
	Status         OrderStatus `json:"status"`
	ChangeDate     time.Time   `json:"change_date"`
}

type OrderItems struct {
	OrderItemID       int     `json:"order_item_id"`
	OrderID           int     `json:"order_id"`
	MenuID            int     `json:"menu_id"`
	Quantity          int     `json:"quantity"`
	PriceAtOrder      float64 `json:"price_at_order"`
	CustomizationInfo string  `json:"customization_info"`
}

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
