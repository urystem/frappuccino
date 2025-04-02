package models

import (
	"github.com/lib/pq"
)

type MenuItem struct {
	ID          uint64               `json:"product_id" db:"id" `
	Name        string               `json:"name" db:"name"`
	Description string               `json:"description" db:"description"`
	Tags        pq.StringArray       `json:"tags" db:"tags"`           /*pgtype.Array[string]*/
	Allergens   pq.StringArray       `json:"allergens" db:"allergens"` /*pgtype.Array[string]*/
	Ingredients []MenuItemIngredient `json:"ingredients"`
	Price       float64              `json:"price" db:"price"`
}

type MenuItemIngredient struct {
	ProductID   uint64  `json:"-" db:"product_id"`
	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
	Quantity    float64 `json:"quantity" db:"quantity"`
}
