package service

import (
	"errors"
	"strings"

	"frappuccino/internal/dal"
	"frappuccino/models"
)

type menuServiceToDal struct {
	menuDal dal.MenuDalInter
}

type MenuServiceInter interface {
	CollectMenus() ([]models.MenuItem, error)
	TakeMenu(uint64) (*models.MenuItem, error)
	DelServiceMenuById(uint64) (*models.MenuDepend, error)
	CreateMenu(*models.MenuItem) ([]models.MenuIngredients, error)
	UpgradeMenu(*models.MenuItem) ([]models.MenuIngredients, error)
}

func ReturnMenuSerStruct(interMenuDal dal.MenuDalInter) MenuServiceInter {
	return &menuServiceToDal{menuDal: interMenuDal}
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

func (ser *menuServiceToDal) CreateMenu(menu *models.MenuItem) ([]models.MenuIngredients, error) {
	ings, err := ser.checkMenuStruct(menu)
	if err != nil {
		return ings, err
	}
	return ser.menuDal.InsertMenu(menu)
}

func (ser *menuServiceToDal) UpgradeMenu(menu *models.MenuItem) ([]models.MenuIngredients, error) {
	ings, err := ser.checkMenuStruct(menu)
	if err != nil {
		return ings, err
	}

	return ser.menuDal.UpdateMenu(menu)
}

func (ser *menuServiceToDal) checkMenuStruct(menu *models.MenuItem) ([]models.MenuIngredients, error) {
	if isInvalidName(menu.Name) {
		return nil, errors.New("invalid name")
	}

	menu.Description = strings.TrimSpace(menu.Description)

	if len(menu.Description) == 0 {
		return nil, errors.New("empty description")
	}
	if len(menu.Tags) == 0 {
		return nil, errors.New("no tags")
	}

	if menu.Price < 0 {
		return nil, errors.New("negative menu price")
	}

	if len(menu.Ingredients) == 0 {
		return nil, errors.New("empty ingridents")
	}

	forTestUniqIngs, invalids := map[uint64]struct{}{}, map[uint64]struct{}{}

	// check for unique and negative quantity ing
	for i, ing := range menu.Ingredients {
		// тазалап алайық, постманнан бар болып келуі мүмкін
		menu.Ingredients[i].Status = nil
		if _, x := forTestUniqIngs[ing.InventoryID]; x || ing.Quantity < 0 {
			if ing.Quantity < 0 {
				menu.Ingredients[i].Status = new(string)
				*menu.Ingredients[i].Status = "invalid quantity"
			}
			invalids[ing.InventoryID] = struct{}{}
		}
		forTestUniqIngs[ing.InventoryID] = struct{}{}
	}
	// if all is correct
	if len(invalids) == 0 {
		return nil, nil
	}

	//"Ленивое удаление" (сдвиг влево)
	// Фильтрация слайса: сдвигаем элементы, которые не нужно удалять, влево
	invalidCount := 0
	for _, ing := range menu.Ingredients {
		// ing := menu.Ingredients[i]
		if _, x := invalids[ing.InventoryID]; x {
			// Если элемент не нужно удалять, ставим его в начало
			if *ing.Status == "" {
				*ing.Status = "Duplicated"
			}
			menu.Ingredients[invalidCount] = ing
			invalidCount++
		}
	}
	return menu.Ingredients[:invalidCount], models.InvalidIngs
}
