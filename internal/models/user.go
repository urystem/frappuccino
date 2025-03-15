package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	UserID     int    `json:"customer_id"`
	Username   string `json:"username"`
	Password   string
	IsAdmin    bool           `json:"is_admin"`
	Age        int            `json:"age"`
	Sex        []uint8        `json:"sex"`
	FirstOrder time.Time      `json:"first_order"`
	Allergens  pq.StringArray `json:"allergens"`
}
