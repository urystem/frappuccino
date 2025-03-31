package dal

import (
	"hot-coffee/models"
)

type OrderDalInter interface {
	WriteOrderDal([]models.Order) error     // Write
	ReadOrdersDal() ([]models.Order, error) // Read
}

func (h *dalCore) ReadOrdersDal() ([]models.Order, error) {
	return nil, nil
}

func (h *dalCore) WriteOrderDal(ords []models.Order) error {
	return nil
}
