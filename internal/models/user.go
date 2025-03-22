package models

import (
	"time"

	"github.com/lib/pq"
)

type Sex int

type User struct {
	UserID     int    `json:"customer_id"`
	Username   string `json:"username"`
	Password   string
	IsAdmin    bool           `json:"is_admin"`
	Age        int            `json:"age"`
	Sex        Sex            `json:"sex"`
	FirstOrder time.Time      `json:"first_order"`
	Allergens  pq.StringArray `json:"allergens"`
}

const (
	male Sex = iota + 1
	female
	undefined
)

func (s Sex) String() string {
	switch s {
	case male:
		return "male"
	case female:
		return "female"
	case undefined:
		return "undefined"
	default:
		return "unknown"
	}
}

func (s Sex) IsValid() bool {
	switch s {
	case male, female, undefined:
		return true
	}
	return false
}
