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
	Created OrderStatus = iota
	Pending
	Processing
	Completed
	Canceled
	Rejected
)
