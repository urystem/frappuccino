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
