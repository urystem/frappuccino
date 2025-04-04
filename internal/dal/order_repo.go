package dal

import (
	"frappuccino/models"
)

type OrderDalInter interface {
	SelectAllOrders() ([]models.Order, error)
	// SelectOrder(uint64) (*models.Order, error)
}

func (core *dalCore) SelectAllOrders() ([]models.Order, error) {
	return nil, nil
}
