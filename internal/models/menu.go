package models

type Menu struct {
	MenuID      int    `json:"menu_id"`
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type Category struct {
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MenuItemsIngredients struct {
	MenuId       int `json:"menu_id"`
	IngredientID int `json:"ingredient_id"`
	Quantity     int `json:"quantity"`
}

type PriceHistory struct {
	PriceHistoryID int `json:"price_history_id"`
	MenuID         int `json:"menu_id"`
	OldPrice       int `json:"old_price"`
	NewPrice       int `json:"new_price"`
	ChangeDate     int `json:"change_date"`
}
