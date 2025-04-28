package models

import (
	"time"

	"github.com/lib/pq"
)

type Order struct {
	ID           uint64         `json:"order_id" db:"id"`                    // Идентификатор заказа
	CustomerName string         `json:"customer_name" db:"customer_name"`    // Имя клиента
	Status       string         `json:"status,omitempty" db:"status"`        // Статус заказа
	Allergens    pq.StringArray `json:"allergens,omitempty" db:"allergens"`  // Список аллергенов
	Total        *float64       `json:"total,omitempty" db:"total"`          // Общая стоимость
	Items        []OrderItem    `json:"items,omitempty"`                     // Заказанные товары (не маппируется на базу)
	CreatedAt    time.Time      `json:"created_at,omitzero" db:"created_at"` // Дата и время создания
	UpdatedAt    time.Time      `json:"updated_at,omitzero" db:"updated_at"` // Дата и время обновления
}

type OrderItem struct {
	Warning       string `json:"error,omitempty"`
	OrderId       uint64 `json:"-" db:"order_id"`
	ProductID     uint64 `json:"product_id" db:"product_id"`
	Quantity      uint64 `json:"quantity" db:"quantity"`
	NotEnoungIngs []struct {
		Inventory_id   uint64  `json:"ingredient_id" db:"id"`
		Inventory_name string  `json:"inventory_name" db:"name"`
		NotEnough      float64 `json:"not_enough" db:"not_enough"`
	} `json:"not_enough,omitempty"`
}

// need not enough inventories
