package models

import "github.com/lib/pq"

type PopularItems struct {
	Items []struct {
		ID    uint64 `json:"item_id" db:"product_id"`
		Name  string `json:"name" db:"name"`
		Count uint64 `json:"count" db:"sum"`
	} `json:"popular_items"`
}

type SearchThings struct {
	Inventories []struct {
		Inventory
		Relevance float64 `json:"relevance" db:"relevance"`
	} `json:"inventory_items,omitempty"`

	Inventory_math *uint64 `json:"inventory_matches,omitempty"`
	Menus          []struct {
		MenuItem
		InventoryItems pq.StringArray `json:"inventories" db:"inventories"`
		Relevance      float64        `json:"relevance" db:"relevance"`
	} `json:"menu_items,omitempty"`

	Menu_math *uint64 `json:"menu_matches,omitempty"`
	Orders    []struct {
		Order
		MenuItems pq.StringArray `json:"menu_items" db:"menu_items"`
		Relevance float64        `json:"relevance" db:"relevance"`
	} `json:"orders,omitempty"`

	Order_math *uint64 `json:"order_matches,omitempty"`
	Total_math uint64  `json:"total_matches"`
}

type OrderStats struct {
	Period     string              `json:"period"`
	Month      string              `json:"month,omitempty"`
	Year       int                 `json:"year,omitempty"`
	OrderItems []map[string]uint64 `json:"ordered_items"`
}

type GetLeftOvers struct {
	SortBy      string `json:"sortedBy"`
	CurrentPage uint64 `json:"currentPage"`
	HasNextPage bool   `json:"hasNextPage"`
	PageSize    uint64 `json:"pageSize"`
	TotalPages  uint64 `json:"totalPages"`
	Data        []struct {
		ID       uint64  `json:"id"  db:"id"`
		Name     string  `json:"name" db:"name"`
		Quantity float64 `json:"quantity" db:"quantity"`
		Price    float64 `json:"price" db:"price"`
	} `json:"data"`
}
