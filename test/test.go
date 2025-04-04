// package main

// import (
// 	"fmt"
// 	"slices"
// )

// type MenuIngredients struct {
// 	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
// 	Quantity    float64 `json:"quantity" db:"quantity"`
// }

// func main() {
// 	// Пример слайса MenuIngredients
// 	ingredients := []MenuIngredients{
// 		{InventoryID: 101, Quantity: 10},
// 		{InventoryID: 102, Quantity: 5},
// 		{InventoryID: 103, Quantity: 7},
// 	}

// 	// Мапа с некорректными inventory_id
// 	invalids := map[uint64]struct{}{
// 		102: {},
// 	}

// 	// Используем slices.DeleteFunc для удаления элементов, которых нет в invalids
// 	ingredients = slices.DeleteFunc(ingredients, func(ing MenuIngredients) bool {
// 		// Удаляем элемент, если его InventoryID нет в invalids
// 		_, exists := invalids[ing.InventoryID]
// 		return !exists
// 	})

// 	// Выводим отфильтрованный слайс
// 	fmt.Println(ingredients)
// }

package main

import "fmt"

type MenuIngredients struct {
	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
	Quantity    float64 `json:"quantity" db:"quantity"`
	Err         string  `json:"err,omitempty"`
}

func main() {
	var str *string
	str = new(string)
	fmt.Println(*str)
	// *str = "ddd"
	// fmt.Println(*str) // выводим значение, на которое указывает str
}
