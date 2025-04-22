package models

type PopularItems struct {
	ID    uint64 `json:"item_id" db:"product_id"`
	Name  string `json:"name" db:"name"`
	Count uint64 `json:"count" db:"sum"`
}
