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

import (
	"fmt"
)

type MenuIngredients struct {
	InventoryID uint64  `json:"inventory_id" db:"inventory_id"`
	Quantity    float64 `json:"quantity" db:"quantity"`
	Err         string  `json:"err,omitempty"`
}

func main() {
	// Пример слайса MenuIngredients
	menu := struct {
		Ingredients []MenuIngredients
	}{
		Ingredients: []MenuIngredients{
			{InventoryID: 101, Quantity: 10},
			{InventoryID: 102, Quantity: 5},
			{InventoryID: 103, Quantity: 7},
			{InventoryID: 104, Quantity: 3},
		},
	}

	// Мапа с некорректными inventory_id
	invalids := map[uint64]struct{}{
		102: {},
	}

	// "Ленивое удаление" с сдвигом влево
	validCount := 0
	for i := range menu.Ingredients {
		ing := menu.Ingredients[i]
		if _, x := invalids[ing.InventoryID]; x {
			if ing.Err == "" {
				ing.Err = "Duplicated"
			}
			menu.Ingredients[validCount] = ing
			validCount++
		}
	}

	// Устанавливаем новую длину слайса
	menu.Ingredients = menu.Ingredients[:validCount]

	// Выводим результат
	fmt.Println(menu.Ingredients)
}
