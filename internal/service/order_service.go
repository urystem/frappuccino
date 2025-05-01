package service

import (
	"fmt"

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
	CreateSomeOrders(batch *models.PostSomeOrders) (*models.OutputBatches, error)
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
	return ser.ordDalInt.InsertOrder(ord, nil)
}

func (ser *ordServiceToDal) UpgradeOrder(id uint64, ord *models.Order) error {
	ord.ID = id
	if err := ser.checkOrderStruct(ord); err != nil {
		return err
	}
	return ser.ordDalInt.UpdateOrder(ord)
}

func (ser *ordServiceToDal) ShutOrder(id uint64) error {
	return ser.ordDalInt.CloseOrder(id)
}

func (ser *ordServiceToDal) CreateSomeOrders(batch *models.PostSomeOrders) (*models.OutputBatches, error) {
	var notValid bool
	for i := range batch.Orders {
		err := ser.checkOrderStruct(&batch.Orders[i])
		if err != nil {
		}
	}
	if notValid {
		return nil, nil
	}

	return ser.ordDalInt.BulkOrderProcessing2(batch.Orders), nil
}

func (ser *ordServiceToDal) checkOrderStruct(ord *models.Order) error {
	if isInvalidName(ord.CustomerName) {
		return fmt.Errorf("%w : invalid name", models.ErrBadInput)
	}
	if len(ord.Items) == 0 {
		return fmt.Errorf("%w : empty items", models.ErrBadInput)
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
	return models.ErrBadInputItems
}
