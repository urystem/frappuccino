package models

import (
	"time"

	"github.com/lib/pq"
)

type Order struct {
	ID           uint64          `json:"order_id" db:"id"`                     // Идентификатор заказа
	CustomerName string          `json:"customer_name" db:"customer_name"`     // Имя клиента
	Status       *string         `json:"status,omitempty" db:"status"`         // Статус заказа
	Allergens    *pq.StringArray `json:"allergens,omitempty" db:"allergens"`   // Список аллергенов
	Total        *float64        `json:"total,omitempty" db:"total"`           // Общая стоимость
	Items        []OrderItem     `json:"items,omitempty"`                      // Заказанные товары (не маппируется на базу)
	CreatedAt    *time.Time      `json:"created_at,omitempty" db:"created_at"` // Дата и время создания
	UpdatedAt    *time.Time      `json:"updated_at,omitempty" db:"updated_at"` // Дата и время обновления
}

type OrderItem struct {
	Err       *string `json:"error,omitempty"`
	ProductID uint64  `json:"product_id" db:"product_id"`
	Quantity  uint64  `json:"quantity" db:"quantity"`
}
