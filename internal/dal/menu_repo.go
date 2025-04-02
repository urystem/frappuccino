package dal

import (
	"fmt"
	"frappuccino/models"
)

type MenuDalInter interface {
	InsertMenu(*models.MenuItem) error       // Write
	SelectMenus() ([]models.MenuItem, error) // Read
}

func (core *dalCore) InsertMenu(menuItems *models.MenuItem) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	insertMenuQ := `
		INSERT INTO menu_items (name, description, tags, allergens, price)
		VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	err = tx.QueryRow(insertMenuQ,
		menuItems.Name,
		menuItems.Description,
		menuItems.Tags,
		menuItems.Allergens,
		menuItems.Price).Scan(&menuItems.ID)
	if err != nil {
		return err
	}
	insert1MenuIngQ := `
		INSERT INTO menu_item_ingredients
		VALUES(:product_id, :inventory_id, :quantity)
	`

	// егер запрос көп болса PrepareNamed дұрыс
	// ал 1 еу ғана бола NamedExec дұрыс
	stmt, err := tx.PrepareNamed(insert1MenuIngQ)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range menuItems.Ingredients {
		v.ProductID = menuItems.ID
		_, err = stmt.Exec(v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (core *dalCore) SelectMenus() ([]models.MenuItem, error) {
	var menus []models.MenuItem
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// var resArray pgtype.Array[string]
	// s := resArray.Elements

	err = tx.Select(&menus, `SELECT * FROM menu_items`)
	if err != nil {
		return nil, err
	}
	MenuIngsQ := `SELECT inventory_id, quantity FROM menu_item_ingredients WHERE product_id=:id`
	stmt, err := tx.PrepareNamed(MenuIngsQ)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := range menus {
		// menus[i] деп структураның өзін бере салдым, (ө)үйткені ол тек 1 ғана аргумент қабылдайды екен
		// menus[i].ID деп query ға $1 қоя салуға келмеді
		stmt.Select(&menus[i].Ingredients, menus[i])
	}
	fmt.Println(menus)
	return menus, tx.Commit()
}
