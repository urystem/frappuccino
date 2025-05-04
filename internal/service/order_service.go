package service

import (
	"errors"
	"fmt"
	"time"

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
	CreateSomeOrders(batch *models.OutputBatches) error
	CollectStatusHistory() ([]models.StatusHistory, error)
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

func (ser *ordServiceToDal) CreateSomeOrders(bulk *models.OutputBatches) error {
	var unknownErr error

	var wasBadInput, wasNotEnough bool

	for i := range bulk.Processed {
		bulk.Summary.TotalOrders++
		bulk.Processed[i].CreatedAt = time.Time{}
		bulk.Processed[i].UpdatedAt = time.Time{}
		err := ser.checkOrderStruct(&bulk.Processed[i])

		if err != nil {
			bulk.Processed[i].Reason = "bad input"
			wasBadInput = true
		} else if err = ser.ordDalInt.InsertOrder(&bulk.Processed[i], &bulk.Summary.InventoryUpdates); err == nil { // err==nil
			bulk.Summary.TotalRevenue += *bulk.Processed[i].Total
			bulk.Summary.Accepted++
			bulk.Processed[i].Status = "accepted"
			bulk.Processed[i].Items = nil
			continue
		} else if errors.Is(err, models.ErrOrderNotEnoughItems) {
			bulk.Processed[i].Reason = "insufficient_inventory"
			wasNotEnough = true
		} else if errors.Is(err, models.ErrNotFoundItems) {
			bulk.Processed[i].Reason = "ErrNotFoundItems"
		} else { // critical error
			unknownErr = err
			bulk.Processed[i].Reason = fmt.Sprintf("unknown error: %s", err.Error())
		}
		bulk.Summary.Rejected++
		bulk.Processed[i].Status = "rejected"

		// bulk.Processed[i].Items = nil
	}
	if unknownErr != nil { // 500
		return unknownErr
	}

	if bulk.Summary.Rejected == 0 { // 200
		return nil
	}

	if bulk.Summary.Accepted != 0 { // 207
		return models.ErrOrdersMultiStatus
	}
	// значить все были Rejected
	if wasBadInput { // 400
		return models.ErrBadInput
	}
	if wasNotEnough {
		return models.ErrOrderNotEnoughItems
	}
	return models.ErrNotFoundItems
}

func (ser *ordServiceToDal) CollectStatusHistory() ([]models.StatusHistory, error) {
	return ser.ordDalInt.SelectAllStatusHistory()
}

func (ser *ordServiceToDal) checkOrderStruct(ord *models.Order) error {
	if isInvalidName(ord.CustomerName) {
		return fmt.Errorf("%w : invalid name", models.ErrBadInput)
	}
	if len(ord.Items) == 0 {
		return fmt.Errorf("%w : empty items", models.ErrBadInput)
	}
	forTestUniqItems := map[uint64]int{}
	var hasZeroQuantity bool
	for i, item := range ord.Items {
		ord.Items[i].Warning = ""
		if item.Quantity == 0 {
			ord.Items[i].Warning = "zero quantity"
			hasZeroQuantity = true
		} else if ind, x := forTestUniqItems[item.ProductID]; x {
			ord.Items[ind].Warning = "duplicated"
			ord.Items[i].Warning = "duplicated"
		}
		forTestUniqItems[item.ProductID] = i
	}

	if len(forTestUniqItems) == len(ord.Items) && !hasZeroQuantity {
		return nil
	}

	var invalids uint64
	for _, item := range ord.Items {
		if len(item.Warning) == 0 {
			ord.Items[invalids] = item
			invalids++
		}
	}
	ord.Items = ord.Items[:invalids]
	return models.ErrBadInputItems
}
