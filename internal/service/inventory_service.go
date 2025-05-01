package service

import (
	"errors"
	"fmt"

	"frappuccino/internal/dal"
	"frappuccino/models"
)

type inventoryServiceDal struct {
	invDal dal.InventoryDataAccess
}

type InventoryService interface {
	CreateInventory(*models.Inventory) error
	CollectInventories() ([]models.Inventory, error)
	TakeInventory(uint64) (*models.Inventory, error)
	UpgradeInventory(*models.Inventory) error
	RemoveInventory(uint64) (*models.InventoryDepend, error)
}

func ReturnInventorySerInt(dalInter dal.InventoryDataAccess) InventoryService {
	return &inventoryServiceDal{invDal: dalInter}
}

func (ser *inventoryServiceDal) CreateInventory(inv *models.Inventory) error {
	if err := ser.checkInventStruct(inv); err != nil {
		return err
	}
	return ser.invDal.InsertInventoryV5(inv)
}

func (ser *inventoryServiceDal) CollectInventories() ([]models.Inventory, error) {
	return ser.invDal.SelectAllInventories()
}

func (ser *inventoryServiceDal) TakeInventory(id uint64) (*models.Inventory, error) {
	inv, err := ser.invDal.SelectInventory(id)
	if errors.Is(err, models.ErrNotFound) {
		err = fmt.Errorf("%w - id = %d", err, id)
	}
	return inv, err
}

func (ser *inventoryServiceDal) UpgradeInventory(inv *models.Inventory) error {
	err := ser.checkInventStruct(inv)
	if err != nil {
		return err
	}
	err = ser.invDal.UpdateInventory(inv)
	if errors.Is(err, models.ErrNotFound) {
		err = fmt.Errorf("%w - id = %d", err, inv.ID)
	}
	return err
}

func (ser *inventoryServiceDal) RemoveInventory(id uint64) (*models.InventoryDepend, error) {
	menuDepend, err := ser.invDal.DeleteInventory(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = fmt.Errorf("%w - id = %d", err, id)
		}
		return nil, err
	}
	if len(menuDepend.Menus) != 0 {
		menuDepend.Err = "Found Depends"
		return menuDepend, nil
	}
	return nil, nil
}

func (ser *inventoryServiceDal) checkInventStruct(inv *models.Inventory) error {
	if isInvalidName(inv.Name) {
		return fmt.Errorf("%w : invalid name - %s", models.ErrBadInput, inv.Name)
	} else if inv.Quantity < 0 {
		return fmt.Errorf("%w : invalid quantity - %f", models.ErrBadInput, inv.Quantity)
	} else if inv.ReorderLvl < 0 {
		return fmt.Errorf("%w : invalid reorder - %f", models.ErrBadInput, inv.ReorderLvl)
	} else if isInvalidName(inv.Unit) {
		return fmt.Errorf("%w : invalid unit - %s", models.ErrBadInput, inv.Unit)
	} else if inv.Price < 0 {
		return fmt.Errorf("%w : invalid price - %f", models.ErrBadInput, inv.Price)
	}
	return nil
}
