package models

import "time"

type Inventory struct {
	ID         uint64  `json:"ingredient_id" db:"id"`
	Name       string  `json:"name" db:"name"`
	Descrip    string  `json:"description" db:"description"`
	Quantity   float64 `json:"quantity" db:"quantity"`
	ReorderLvl float64 `json:"reorder_level" db:"reorder_level"`
	Unit       string  `json:"unit" db:"unit"`
	Price      float64 `json:"price" db:"price"`
}

// бұғаен json тегі қатты керек емес)
type InventoryTransaction struct {
	ID             uint64    `db:"id" json:"id"`
	InventoryID    uint64    `db:"inventory_id" json:"inventory_id"`
	QuantityChange float64   `db:"quantity_change" json:"quantity_change"`
	Reason         string    `db:"reason" json:"reason"` // ENUM в БД, в Go — string
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type InventoryDepend struct {
	Err   string `json:"error"`
	Menus []struct {
		ProductID uint64 `json:"product_id" db:"id"`
		Name      string `json:"name" db:"name"`
	} `json:"menu_items"`
}
