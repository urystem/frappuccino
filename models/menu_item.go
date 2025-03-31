package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MenuItem struct {
	ID          uint64               `json:"product_id" db:"id" `
	Name        string               `json:"name" db:"name"`
	Description string               `json:"description" db:"description"`
	Tags        pgtype.Array[string] `json:"tags" db:"tags"`
	Allergens   []string             `json:"allergens" db:"allergens"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
	Price       float64              `json:"price" db:"price"`
}

type MenuItemIngredient struct {
	ProductID   uint64  `json:"-" db:"product_id"`
	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
	Quantity    float64 `json:"quantity" db:"quantity"`
}
