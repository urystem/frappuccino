package models

import (
	"errors"
)

var (
	ErrDelDepend       = errors.New("cannot delete")
	ErrConflict        = errors.New(" already exists")
	ErrNotFound        = errors.New(" not found")
	ErrIngsNotFound    = errors.New("ing not found")
	InvalidIngs        = errors.New(" not found ingrident(s) for menu")
	ErrOrdNotFoundItem = errors.New(" not found item in menu")
	ErrOrdStatusClosed = errors.New(" closed")
	ErrOrdNotEnough    = errors.New(" Not enough quantity")
	ErrOrderItems      = errors.New("was error")
	ErrInputOrder      = errors.New("invalid input order struct")
)
