package models

import "errors"

var (
	ErrConflict        = errors.New(" already exists")
	ErrNotFound        = errors.New(" not found")
	ErrNotFoundIngs    = errors.New(" not found ingrident(s) for menu")
	ErrOrdNotFoundItem = errors.New(" not found item in menu")
	ErrOrdStatusClosed = errors.New(" closed")
	ErrOrdNotEnough    = errors.New(" Not enough quantity")
)

// "d  ff a the order status is open" //32-24=8
//
