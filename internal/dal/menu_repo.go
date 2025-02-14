package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

type menuDalStruct struct {
	menuFilePath string
}

type MenuDalInter interface {
	WriteMenuDal([]models.MenuItem) error    // Write
	ReadMenuDal() ([]models.MenuItem, error) // Read
}

func ReturnMenuDalStruct(filepath string) *menuDalStruct {
	return &menuDalStruct{menuFilePath: filepath}
}

func (md *menuDalStruct) WriteMenuDal(menuItems []models.MenuItem) error {
	file, err := os.Create(md.menuFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(menuItems)
}

func (md *menuDalStruct) ReadMenuDal() ([]models.MenuItem, error) {
	file, err := os.Open(md.menuFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var menuItems []models.MenuItem
	if err = json.NewDecoder(file).Decode(&menuItems); err != nil {
		return nil, err
	}
	return menuItems, nil
}
