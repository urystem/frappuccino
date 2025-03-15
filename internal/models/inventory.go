package models

import "github.com/lib/pq"

type Inventory struct {
	InventoryId   int            `json:"inventory_id"`
	InventoryName string         `json:"inventory_name"`
	Quantity      float64        `json:"quantity"`
	Unit          string         `json:"unit"`
	Allergens     pq.StringArray `json:"allergens"`
	IsActive      bool           `json:"is_active"`
}

type InventoryTransactions struct {
	TransactionId string          `json:"transaction_id"`
	InventoryId   int             `json:"inventory_id"`
	Type          TransactionType `json:"type"`
	Quantity      float64         `json:"quantity"`
	Date          string          `json:"date"`
}

type TransactionType int

const (
	Subtract TransactionType = iota
	Add
	Zero
)
