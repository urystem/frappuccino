package models

import (
	"errors"
	"time"

	"github.com/lib/pq"
)

type InventoryItem struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Quantity  int            `json:"quantity"`
	Unit      string         `json:"unit"`
	Allergens pq.StringArray `json:"allergens" db:"allergens"`
	ExtraInfo JSONB          `json:"extra_info"`
}

type InventoryTransaction struct {
	ID              int       `json:"id"`
	InventoryItemID int       `json:"inventory_item_id"`
	QuantityChange  int       `json:"quantity_change"`
	TransactionDate time.Time `json:"transaction_date"`
}

func (i *InventoryItem) IsValid() error {
	if i.Name == "" || i.Unit == "" {
		return errors.New("name or unit cannot be empty")
	}

	if i.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return nil
}
