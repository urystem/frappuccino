package models

import "github.com/lib/pq"

// type TotalOrders struct {
// 	Total float64
// }

type PopularItems struct {
	Items []struct {
		ID    uint64 `json:"item_id" db:"product_id"`
		Name  string `json:"name" db:"name"`
		Count uint64 `json:"count" db:"sum"`
	} `json:"popular_items"`
}

// type CountOfOrderedItem map[string]uint64

type SearchThings struct {
	Inventories []struct {
		Inventory
		Relevance float64 `json:"relevance" db:"relevance"`
	} `json:"inventory_items,omitempty"`

	Menus []struct {
		MenuItem
		InventoryItems pq.StringArray `json:"inventories" db:"inventories"`
		Relevance      float64        `json:"relevance" db:"relevance"`
	} `json:"menu_items,omitempty"`

	// Orders     []OrderSearch `json:"orders"`
	Total_math uint64 `json:"total_matches"`
}
