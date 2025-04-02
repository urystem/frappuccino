package service

import (
	"errors"

	"frappuccino/internal/dal"
	"frappuccino/models"
)

type menuServiceToDal struct {
	menuDal dal.MenuDalInter
}

type MenuServiceInter interface {
	CreateMenu(*models.MenuItem) ([]uint64, error)
	CollectMenus() ([]models.MenuItem, error)
	TakeMenu(uint64) (*models.MenuItem, error)
	DelServiceMenuById(uint64) (*models.MenuDepend, error)
	// PutServiceMenuById(*models.MenuItem, string) ([]string, error)
}

func ReturnMenuSerStruct(interMenuDal dal.MenuDalInter) *menuServiceToDal {
	return &menuServiceToDal{menuDal: interMenuDal}
}

func (ser *menuServiceToDal) CreateMenu(menu *models.MenuItem) ([]uint64, error) {
	if isInvalidName(menu.Name){
		return nil, errors.New("ff")
	}
	return ser.menuDal.InsertMenu(menu)
}

func (ser *menuServiceToDal) CollectMenus() ([]models.MenuItem, error) {
	return ser.menuDal.SelectAllMenus()
}

func (ser *menuServiceToDal) TakeMenu(id uint64) (*models.MenuItem, error) {
	return ser.menuDal.SelectMenu(id)
}

func (ser *menuServiceToDal) DelServiceMenuById(id uint64) (*models.MenuDepend, error) {
	return ser.menuDal.DeleteMenu(id)
}

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

func checkMenuStruct(menu models.MenuItem) error {
	if isInvalidName(menu.Name) {
		return errors.New("invalid name")
	}
	if 
}
