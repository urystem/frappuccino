package service

import (
	"errors"

	"frappuccino/internal/dal"
	"frappuccino/models"
)

type inventoryServiceDal struct {
	invDal dal.InventoryDataAccess
}

func ReturnInventorySerStruct(dalInter dal.InventoryDataAccess) *inventoryServiceDal {
	return &inventoryServiceDal{invDal: dalInter}
}

type InventoryService interface {
	CreateInventory(*models.Inventory) error
	CollectInventories() ([]models.Inventory, error)
	TakeInventory(uint64) (*models.Inventory, error)
	UpgradeInventory(*models.Inventory) error
	RemoveInventory(uint64) (*models.InventoryDepend, error)
	// PutAllInvets([]models.InventoryItem) ([]string, error)
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
	return ser.invDal.SelectInventory(id)
}

func (ser *inventoryServiceDal) UpgradeInventory(inv *models.Inventory) error {
	if err := ser.checkInventStruct(inv); err != nil {
		return err
	}
	return ser.invDal.UpdateInventory(inv)
}

func (ser *inventoryServiceDal) RemoveInventory(id uint64) (*models.InventoryDepend, error) {
	return ser.invDal.DeleteInventory(id)
}

// func (ser *inventoryServiceDal) PutAllInvets(invents []models.InventoryItem) ([]string, error) {
// 	items, err := ser.dal.ReadInventory()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var notFounds []string
// 	for _, invent := range invents {
// 		var isHere bool
// 		for _, item := range items {
// 			if invent.IngredientID == item.IngredientID {
// 				isHere = true
// 				break
// 			}
// 		}
// 		if !isHere {
// 			notFounds = append(notFounds, invent.IngredientID)
// 		}
// 	}
// 	if len(notFounds) != 0 {
// 		return notFounds, models.ErrNotFoundIngs
// 	}
// 	for _, invent := range invents {
// 		if err := ser.UpdateInventory(invent.IngredientID, &invent); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return nil, nil
// }

func (ser *inventoryServiceDal) checkInventStruct(inv *models.Inventory) error {
	if isInvalidName(inv.Name) {
		return errors.New("invalid name")
	} else if inv.Quantity < 0 {
		return errors.New("invalid quantity")
	} else if inv.ReorderLvl < 0 {
		return errors.New("invalid reorder")
	} else if isInvalidName(inv.Unit) {
		return errors.New("invalid unit")
	} else if inv.Price < 0 {
		return errors.New("invalid price")
	}
	return nil
}
