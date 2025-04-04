package models

import "github.com/lib/pq"

type Order struct {
	ID           string `json:"order_id"`
	CustomerName string `json:"customer_name"`
	Status       string `json:"status"`
	Allergens    pq.StringArray
	// total 
	Items        []OrderItem `json:"items"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
