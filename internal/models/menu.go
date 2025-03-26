package models

import (
	"time"

	"github.com/lib/pq"
)

type MenuItem struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Allergens   pq.StringArray       `json:"allergens"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}

type MenuItemIngredient struct {
	ID              int `json:"id"`
	MenuItemID      int `json:"menu_item_id"`
	InventoryItemID int `json:"inventory_item_id"`
	Quantity        int `json:"quantity"`
}

type PriceHistory struct {
	ID         int       `json:"id"`
	MenuItemID int       `json:"menu_item_id"`
	OldPrice   float64   `json:"old_price"`
	NewPrice   float64   `json:"new_price"`
	ChangedAt  time.Time `json:"changed_at"`
}
