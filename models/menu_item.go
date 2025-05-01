package models

import (
	"github.com/lib/pq"
)

type MenuItem struct {
	ID          uint64            `json:"product_id" db:"id" `
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Tags        pq.StringArray    `json:"tags" db:"tags"`           /*pgtype.Array[string]*/
	Allergens   pq.StringArray    `json:"allergens" db:"allergens"` /*pgtype.Array[string]*/
	Price       float64           `json:"price" db:"price"`
	Ingredients []MenuIngredients `json:"ingredients,omitempty"`
}

type MenuIngredients struct {
	Status      string  `json:"status,omitempty"` // егер нил болса мүлдем жасырып тастайды
	ProductID   uint64  `json:"-" db:"product_id"`
	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
	Quantity    float64 `json:"quantity" db:"quantity"`
}

type MenuDepend struct {
	Err    string  `json:"error"`
	Orders []Order `json:"orders"`
}

// var (
// 	ErrMenuInput  = errors.New("bad request")
// 	ErrMenuName   = fmt.Errorf("%w: invalid name", ErrMenuInput)
// 	ErrMenuDesc   = fmt.Errorf("%w: invalid description", ErrMenuInput)
// 	ErrMenuNoTags = fmt.Errorf("%w: no tags", ErrMenuInput)
// 	ErrMenuPrice  = fmt.Errorf("%w: negative menu price", ErrMenuInput)
// 	ErrMenuIngs   = fmt.Errorf("%w: empty ingridents", ErrMenuInput)

// 	ErrMenuNameConflict = errors.New("conflict")
// 	ErrMenuNotFound     = errors.New("menu not found")
// )
