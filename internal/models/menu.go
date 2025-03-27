package models

import (
	"errors"
	"time"

	"github.com/lib/pq"
)

type MenuItem struct {
	ID          int                  `json:"id" db:"menu_item_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Details     JSONB                `json:"details"`
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
	OldPrice   float32   `json:"old_price"`
	NewPrice   float32   `json:"new_price"`
	ChangedAt  time.Time `json:"changed_at"`
}

func (m *MenuItem) IsValid() error {
	if m.Name == "" || m.Description == "" {
		return errors.New("name or description cannot be empty")
	}

	if m.Price <= 0 {
		return errors.New("price cannot be less than or equal to zero")
	}

	if len(m.Ingredients) == 0 {
		return errors.New("ingridients cannot be empty")
	}

	return nil
}
