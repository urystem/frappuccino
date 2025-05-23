package service

import (
	"fmt"
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
	CreateMenu(*models.MenuItem) error
	UpgradeMenu(*models.MenuItem) error
	CollectHistory() ([]models.PriceHistory, error)
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

func (ser *menuServiceToDal) CreateMenu(menu *models.MenuItem) error {
	err := ser.checkMenuStruct(menu)
	if err != nil {
		return err
	}
	return ser.menuDal.InsertMenu(menu)
}

func (ser *menuServiceToDal) UpgradeMenu(menu *models.MenuItem) error {
	err := ser.checkMenuStruct(menu)
	if err != nil {
		return err
	}

	return ser.menuDal.UpdateMenu(menu)
}

func (ser *menuServiceToDal) CollectHistory() ([]models.PriceHistory, error) {
	return ser.menuDal.SelectPriceHistory()
}

func (ser *menuServiceToDal) checkMenuStruct(menu *models.MenuItem) error {
	if isInvalidName(menu.Name) {
		return fmt.Errorf("%w: invalid name - %s", models.ErrBadInput, menu.Name)
	}

	menu.Description = strings.TrimSpace(menu.Description)

	if len(menu.Description) == 0 {
		return fmt.Errorf("%w: empty description", models.ErrBadInput)
	}
	if len(menu.Tags) == 0 {
		return fmt.Errorf("%w: no tags", models.ErrBadInput)
	}

	if menu.Price < 0 {
		return fmt.Errorf("%w: negative menu price", models.ErrBadInput)
	}

	if len(menu.Ingredients) == 0 {
		return fmt.Errorf("%w: empty ingridents", models.ErrBadInput)
	}

	forTestUniqIngs, invalids := map[uint64]struct{}{}, map[uint64]struct{}{}

	// check for unique and negative quantity ing
	for i, ing := range menu.Ingredients {
		// тазалап алайық, постманнан бар болып келуі мүмкін
		menu.Ingredients[i].Status = ""
		if _, x := forTestUniqIngs[ing.InventoryID]; x || ing.Quantity < 0 {
			if ing.Quantity < 0 {
				menu.Ingredients[i].Status = "invalid quantity"
			}
			invalids[ing.InventoryID] = struct{}{}
		}
		forTestUniqIngs[ing.InventoryID] = struct{}{}
	}
	// if all is correct
	if len(invalids) == 0 {
		return nil
	}

	//"Ленивое удаление" (сдвиг влево)
	// Фильтрация слайса: сдвигаем элементы, которые не нужно удалять, влево
	var invalidCount uint64
	for _, ing := range menu.Ingredients {
		// ing := menu.Ingredients[i]
		if _, x := invalids[ing.InventoryID]; x {
			// Если элемент не нужно удалять, ставим его в начало
			if ing.Status == "" {
				ing.Status = "Duplicated"
			}
			menu.Ingredients[invalidCount] = ing
			invalidCount++
		}
	}
	menu.Ingredients = menu.Ingredients[:invalidCount]
	return models.ErrBadInputItems
}
