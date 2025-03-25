package models

type Inventory struct {
	InventoryId   int     `json:"inventory_id"`
	InventoryName string  `json:"inventory_name"`
	Quantity      float64 `json:"quantity"`
	Unit          string  `json:"unit"`
}

type InventoryTransactions struct {
	TransactionId string  `json:"transaction_id"`
	InventoryId   int     `json:"inventory_id"`
	Quantity      float64 `json:"quantity"`
	Date          string  `json:"date"`
}
