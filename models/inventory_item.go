package models

type InventoryItem struct {
	ID         uint    `json:"ingredient_id"`
	Name       string  `json:"name"`
	Descrip    string  `json:"description"`
	Quantity   float64 `json:"quantity"`
	ReorderLvl float64 `json:"reorder_level"`
	Unit       string  `json:"unit"`
	Price      float64 `json:"price"`
}
