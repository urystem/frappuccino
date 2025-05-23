package models

import (
	"time"

	"github.com/lib/pq"
)

type Order struct {
	ID           uint64         `json:"order_id" db:"id"`                   // Идентификатор заказа
	CustomerName string         `json:"customer_name" db:"customer_name"`   // Имя клиента
	Status       string         `json:"status,omitempty" db:"status"`       // Статус заказа
	Allergens    pq.StringArray `json:"allergens,omitempty" db:"allergens"` // Список аллергенов
	Reason       string         `json:"reason,omitempty"`
	Total        *float64       `json:"total,omitempty" db:"total"`          // Общая стоимость
	Items        []OrderItem    `json:"items,omitempty"`                     // Заказанные товары (не маппируется на базу)
	CreatedAt    time.Time      `json:"created_at,omitzero" db:"created_at"` // Дата и время создания
	UpdatedAt    time.Time      `json:"updated_at,omitzero" db:"updated_at"` // Дата и время обновления
}

type OrderItem struct {
	Warning       string         `json:"error,omitempty"`
	// OrderId       uint64         `json:"-" db:"order_id"`
	ProductID     uint64         `json:"product_id" db:"product_id"`
	Quantity      uint64         `json:"quantity,omitempty" db:"quantity"`
	Allergens     pq.StringArray `json:"allergens,omitempty" db:"allergens"`
	NotEnoungIngs []struct {     // қарау керек: егер 2 orderItem де бірдей Inventory болса жетіспейтіндері NotEnough әртүрлі болады
		Inventory_id   uint64  `json:"ingredient_id" db:"id"`
		Inventory_name string  `json:"inventory_name" db:"name"`
		NotEnough      float64 `json:"not_enough" db:"not_enough"`
	} `json:"not_enough,omitempty"`
}

// input
type PostSomeOrders struct {
	Orders []Order `json:"orders"`
}

type InventoryUpdate struct {
	InventoryID  uint64  `json:"ingredient_id" db:"id"`
	Name         string  `json:"name" db:"name"`
	QuantityUsed float64 `json:"quantity_used" db:"quantity_used"`
	Remaining    float64 `json:"remaining" db:"remaining"`
}

// output
type OutputBatches struct {
	Processed []Order `json:"processed_orders"`

	Summary struct {
		TotalOrders      uint64            `json:"total_orders"`
		Accepted         uint64            `json:"accepted"`
		Rejected         uint64            `json:"rejected"`
		TotalRevenue     float64           `json:"total_revenue"`
		InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
	} `json:"summary"`
}

type StatusHistory struct {
	ID      uint64    `json:"history_id" db:"id"`
	OrderID uint64    `json:"order_id" db:"order_id"`
	Status  string    `json:"status" db:"status"`
	Updated time.Time `json:"updated_at" db:"updated_at"`
}
