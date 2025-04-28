package service

import (
	"errors"

	"frappuccino/internal/dal"
	"frappuccino/models"
)

type ordServiceToDal struct {
	ordDalInt dal.OrderDalInter
}

type OrdServiceInter interface {
	CollectOrders() ([]models.Order, error)
	TakeOrder(uint64) (*models.Order, error)
	RemoveOrder(uint64) error
	CreateOrder(*models.Order) error
	UpgradeOrder(id uint64, ord *models.Order) error
	ShutOrder(uint64) error
}

func ReturnOrdSerStruct(ord dal.OrderDalInter) OrdServiceInter {
	return &ordServiceToDal{ordDalInt: ord}
}

func (ser *ordServiceToDal) CollectOrders() ([]models.Order, error) {
	return ser.ordDalInt.SelectAllOrders()
}

func (ser *ordServiceToDal) TakeOrder(id uint64) (*models.Order, error) {
	return ser.ordDalInt.SelectOrder(id)
}

func (ser *ordServiceToDal) RemoveOrder(id uint64) error {
	return ser.ordDalInt.DeleteOrder(id)
}

func (ser *ordServiceToDal) CreateOrder(ord *models.Order) error {
	if err := ser.checkOrderStruct(ord); err != nil {
		return err
	}
	return ser.ordDalInt.InsertOrder(ord)
}

func (ser *ordServiceToDal) UpgradeOrder(id uint64, ord *models.Order) error {
	if err := ser.checkOrderStruct(ord); err != nil {
		return err
	}
	return ser.ordDalInt.UpdateOrder(id, ord)
}

func (ser *ordServiceToDal) ShutOrder(id uint64) error {
	return ser.ordDalInt.CloseOrder(id)
}

func (ser *ordServiceToDal) checkOrderStruct(ord *models.Order) error {
	if isInvalidName(ord.CustomerName) {
		return errors.New("invalid name")
	}
	if len(ord.Items) == 0 {
		return errors.New("emyty items")
	}
	forTestUniqItems := map[uint64]int{}
	for i, item := range ord.Items {
		ord.Items[i].Warning = ""
		if ind, x := forTestUniqItems[item.ProductID]; x {
			ord.Items[ind].Warning = "duplicated"
			ord.Items[i].Warning = "duplicated"
		}
		forTestUniqItems[item.ProductID] = i
	}

	if len(forTestUniqItems) == len(ord.Items) {
		return nil
	}

	var duplicateds uint64
	for _, item := range ord.Items {
		if _, x := forTestUniqItems[item.ProductID]; x {
			ord.Items[duplicateds] = item
			duplicateds++
		}
	}
	ord.Items = ord.Items[:duplicateds]
	return models.ErrInputOrder
}
