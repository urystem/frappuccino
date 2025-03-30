package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type menuServiceToDal struct {
	menuDalInt   dal.MenuDalInter
	inventDalInt dal.InventoryDataAccess
}

type MenuServiceInter interface {
	PostServiceMenu(*models.MenuItem) ([]string, error)
	GetServiceMenus() ([]models.MenuItem, error)
	GetServiceMenuById(string) (*models.MenuItem, error)
	PutServiceMenuById(*models.MenuItem, string) ([]string, error)
	DelServiceMenuById(string) error
}

func ReturnMenuSerStruct(interMenuDal dal.MenuDalInter, interIngDal dal.InventoryDataAccess) *menuServiceToDal {
	return &menuServiceToDal{menuDalInt: interMenuDal, inventDalInt: interIngDal}
}

// func (ser *menuServiceToDal) PostServiceMenu(menu *models.MenuItem) ([]string, error) {
// 	if ings, err := ser.checkNotFoundIngs(menu.Ingredients); err != nil {
// 		return nil, err
// 	} else if len(ings) != 0 {
// 		return ings, models.ErrNotFoundIngs
// 	}
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, v := range menus {
// 		if v.ID == menu.ID {
// 			return nil, models.ErrConflict
// 		}
// 	}
// 	return nil, ser.menuDalInt.WriteMenuDal(append(menus, *menu))
// }

// func (ser *menuServiceToDal) GetServiceMenus() ([]models.MenuItem, error) {
// 	return ser.menuDalInt.ReadMenuDal()
// }

// func (ser *menuServiceToDal) GetServiceMenuById(id string) (*models.MenuItem, error) {
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, v := range menus {
// 		if v.ID == id {
// 			return &v, nil
// 		}
// 	}
// 	return nil, models.ErrNotFound
// }

// func (ser *menuServiceToDal) PutServiceMenuById(menu *models.MenuItem, id string) ([]string, error) {
// 	if ings, err := ser.checkNotFoundIngs(menu.Ingredients); err != nil {
// 		return nil, err
// 	} else if ings != nil && len(ings) != 0 {
// 		return ings, models.ErrNotFoundIngs
// 	}
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i, v := range menus {
// 		if v.ID == id {
// 			menu.ID = id
// 			menus[i] = *menu
// 			return nil, ser.menuDalInt.WriteMenuDal(menus)
// 		}
// 	}
// 	return nil, models.ErrNotFound
// }

// func (ser *menuServiceToDal) DelServiceMenuById(id string) error {
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return err
// 	}
// 	for i, v := range menus {
// 		if v.ID == id {
// 			return ser.menuDalInt.WriteMenuDal(append(menus[:i], menus[i+1:]...))
// 		}
// 	}
// 	return models.ErrNotFound
// }

// func (ser *menuServiceToDal) checkNotFoundIngs(itemsToCheck []models.MenuItemIngredient) ([]string, error) {
// 	ingDul, err := ser.inventDalInt.ReadInventory()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var notFoundIngs []string
// 	for _, ing := range itemsToCheck {
// 		var isHere bool
// 		for _, ingInDal := range ingDul {
// 			if ing.IngredientID == ingInDal.IngredientID {
// 				isHere = true
// 				break
// 			}
// 		}
// 		if !isHere {
// 			notFoundIngs = append(notFoundIngs, ing.IngredientID)
// 		}
// 	}
// 	return notFoundIngs, nil
// }
