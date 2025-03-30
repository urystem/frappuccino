package models

import (
	"errors"
)

var (
	ErrInvalidName     = errors.New("")
	ErrConflict        = errors.New(" already exists")
	ErrNotFound        = errors.New(" not found")
	ErrNotFoundIngs    = errors.New(" not found ingrident(s) for menu")
	ErrOrdNotFoundItem = errors.New(" not found item in menu")
	ErrOrdStatusClosed = errors.New(" closed")
	ErrOrdNotEnough    = errors.New(" Not enough quantity")
)
