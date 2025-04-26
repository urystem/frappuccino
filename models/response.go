package models

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
		Relevance float64 `db:"relevance"`
	} `json:"inventory_items"`
	// Menus      []MenuSearch  `json:"menu_items"`
	// Orders     []OrderSearch `json:"orders"`
	Total_math uint64 `json:"total_matches"`
}

type MenuSearch struct {
	ID        uint64
	Name      string
	Desc      string
	Price     float64
	Relevance float64
}

type OrderSearch struct {
	ID            uint64
	Customer_name string
	Items         []string
	Total         float64
	Relevance     float64
}
