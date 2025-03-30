package service

import (
	"errors"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryService interface {
	InsertInventory(*models.InventoryItem) error
	// GetAllInventory() ([]models.InventoryItem, error)
	// GetSpecificInventory(string) (*models.InventoryItem, error)
	// UpdateInventory(string, *models.InventoryItem) error
	// DeleteInventory(string) error
	// PutAllInvets([]models.InventoryItem) ([]string, error)
}

type inventoryServiceDal struct {
	dal dal.InventoryDataAccess
}

func NewInventoryService(dalInter dal.InventoryDataAccess) *inventoryServiceDal {
	return &inventoryServiceDal{dal: dalInter}
}

func (ser *inventoryServiceDal) InsertInventory(inv *models.InventoryItem) error {
	if err := checkInventStruct(inv); err != nil {
		return err
	}
	_, err := ser.dal.InsertInventory(inv)
	return err
}

// func (ser *inventoryServiceDal) GetAllInventory() ([]models.InventoryItem, error) {
// 	return ser.dal.ReadInventory()
// }

// func (ser *inventoryServiceDal) GetSpecificInventory(id string) (*models.InventoryItem, error) {
// 	invents, err := ser.dal.ReadInventory()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, v := range invents {
// 		if v.IngredientID == id {
// 			return &v, nil
// 		}
// 	}
// 	return nil, models.ErrNotFound
// }

// func (ser *inventoryServiceDal) UpdateInventory(id string, inv *models.InventoryItem) error {
// 	items, err := ser.dal.ReadInventory()
// 	if err != nil {
// 		return err
// 	}
// 	for i, v := range items {
// 		if v.IngredientID == id {
// 			inv.IngredientID = id
// 			items[i] = *inv
// 			return ser.dal.WriteInventory(items)
// 		}
// 	}
// 	return models.ErrNotFound
// }

// func (ser *inventoryServiceDal) DeleteInventory(id string) error {
// 	items, err := ser.dal.ReadInventory()
// 	if err != nil {
// 		return err
// 	}
// 	for i, v := range items {
// 		if v.IngredientID == id {
// 			return ser.dal.WriteInventory(append(items[:i], items[i+1:]...))
// 		}
// 	}
// 	return models.ErrNotFound
// }

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

func checkInventStruct(inv *models.InventoryItem) error {
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
