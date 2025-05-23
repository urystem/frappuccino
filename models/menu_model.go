package models

import (
	"time"

	"github.com/lib/pq"
)

type MenuItem struct {
	ID          uint64            `json:"product_id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Tags        pq.StringArray    `json:"tags" db:"tags"`           /*pgtype.Array[string]*/
	Allergens   pq.StringArray    `json:"allergens" db:"allergens"` /*pgtype.Array[string]*/
	Price       float64           `json:"price" db:"price"`
	Ingredients []MenuIngredients `json:"ingredients,omitempty"`
}

type MenuIngredients struct {
	Status      string  `json:"status,omitempty"` // егер нил болса мүлдем жасырып тастайды
	ProductID   uint64  `json:"-" db:"product_id"`
	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
	Quantity    float64 `json:"quantity" db:"quantity"`
}

type MenuDepend struct {
	Err    string  `json:"error"`
	Orders []Order `json:"orders"`
}

type PriceHistory struct {
	ID        uint64    `json:"history_id" db:"id"`
	ProductID uint64    `json:"product_id" db:"product_id"`
	OldPrice  float64   `json:"old_price" db:"old_price"`
	NewPrice  float64   `json:"new_price" db:"new_price"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
