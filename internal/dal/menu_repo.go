package dal

import (
	"frappuccino/models"
)

type MenuDalInter interface {
	InsertMenu(*models.MenuItem) error
	SelectAllMenus() ([]models.MenuItem, error)
	SelectMenu(uint64) (*models.MenuItem, error)
	DeleteMenu(uint64) (*models.MenuDepend, error)
}

func (core *dalCore) InsertMenu(menuItems *models.MenuItem) ([]uint64, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	insertMenuQ := `
		INSERT INTO menu_items (name, description, tags, allergens, price)
		VALUES ($1, $2, $3, $4, $5)
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
		VALUES(:product_id, :inventory_id, :quantity)`

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

func (core *dalCore) SelectAllMenus() ([]models.MenuItem, error) {
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

	stmt, err := tx.PrepareNamed(`SELECT inventory_id, quantity FROM menu_item_ingredients WHERE product_id=:id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := range menus {
		// menus[i] деп структураның өзін бере салдым, (ө)үйткені ол тек 1 ғана аргумент қабылдайды екен
		// menus[i].ID деп query ға $1 қоя салуға келмеді
		stmt.Select(&menus[i].Ingredients, menus[i])
	}
	return menus, tx.Commit()
}

func (core *dalCore) SelectMenu(id uint64) (*models.MenuItem, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var menu models.MenuItem

	err = tx.Get(&menu, `SELECT * FROM menu_items WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	err = tx.Select(&menu.Ingredients, `SELECT * FROM menu_item_ingredients WHERE product_id=$1`, id)
	if err != nil {
		return nil, err
	}
	return &menu, tx.Commit()
}

func (core *dalCore) DeleteMenu(id uint64) (*models.MenuDepend, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	SELECT order_id, customer_name 
		FROM order_items 
		JOIN orders ON order_id=id 
		WHERE status <> 'processing' AND product_id=$1`

	var menuDepend models.MenuDepend
	err = tx.Select(&menuDepend.Orders, query, id)
	if err != nil {
		return nil, err
	}
	if len(menuDepend.Orders) != 0 {
		menuDepend.Err = "found depends"
		return &menuDepend, nil
	}

	res, err := tx.Exec(`DELETE from menu_items WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	affects, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affects == 0 {
		return nil, models.ErrNotFound
	}

	return nil, tx.Commit()
}
