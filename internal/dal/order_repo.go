package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

type orderDalStruct struct {
	ordersFilePath string
}

type OrderDalInter interface {
	WriteOrderDal([]models.Order) error     // Write
	ReadOrdersDal() ([]models.Order, error) // Read
}

func ReturnOrdDalStruct(filepath string) *orderDalStruct {
	return &orderDalStruct{ordersFilePath: filepath}
}

func (h *orderDalStruct) ReadOrdersDal() ([]models.Order, error) {
	file, err := os.Open(h.ordersFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var orders []models.Order
	if err = json.NewDecoder(file).Decode(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (h *orderDalStruct) WriteOrderDal(ords []models.Order) error {
	file, err := os.Create(h.ordersFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(ords)
}
