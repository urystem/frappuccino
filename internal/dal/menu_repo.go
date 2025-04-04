package dal

import (
	"database/sql"
	"fmt"

	"frappuccino/models"

	"github.com/jmoiron/sqlx"
)

type MenuDalInter interface {
	SelectAllMenus() ([]models.MenuItem, error)
	SelectMenu(uint64) (*models.MenuItem, error)
	DeleteMenu(uint64) (*models.MenuDepend, error)
	InsertMenu(*models.MenuItem) ([]models.MenuIngredients, error)
	UpdateMenu(*models.MenuItem) ([]models.MenuIngredients, error)
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
		err = stmt.Select(&menus[i].Ingredients, menus[i])
		if err != nil {
			return nil, nil
		}
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

func (core *dalCore) InsertMenu(menuItems *models.MenuItem) ([]models.MenuIngredients, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	notFoundIngs, err := core.checkIngs(tx, menuItems.Ingredients)
	if err != nil {
		return notFoundIngs, err
	}

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
		return nil, err
	}
	err = core.insertToMenuIngs(tx, menuItems.ID, menuItems.Ingredients)
	if err != nil {
		return nil, err
	}
	return nil, tx.Commit()
}

func (core *dalCore) UpdateMenu(menuItems *models.MenuItem) ([]models.MenuIngredients, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	notFoundIngs, err := core.checkIngs(tx, menuItems.Ingredients)
	if err != nil {
		return notFoundIngs, err
	}

	updateMenuQ := `
	UPDATE menu_items 
		SET name=:name, description = :description, 
			tags = :tags, allergens = :allergens, price = :price
		WHERE id = :id`

	result, err := tx.NamedExec(updateMenuQ, menuItems)
	if err != nil {
		fmt.Println("qalaisyn")
		return nil, err
	}

	affects, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affects == 0 {
		return nil, models.ErrNotFound
	}

	_, err = tx.Exec(`
	DELETE FROM menu_item_ingredients
		WHERE product_id = $1
	`, menuItems.ID)
	if err != nil {
		return nil, err
	}
	err = core.insertToMenuIngs(tx, menuItems.ID, menuItems.Ingredients)
	if err != nil {
		return nil, err
	}
	return nil, tx.Commit()
}

func (core *dalCore) checkIngs(tx *sqlx.Tx, ings []models.MenuIngredients) ([]models.MenuIngredients, error) {
	stmt, err := tx.Prepare(`SELECT TRUE FROM inventory WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var notFoundCount uint64
	for _, v := range ings {
		var exists bool

		err = stmt.QueryRow(v.InventoryID).Scan(&exists)

		if err == sql.ErrNoRows {
			*v.Status = "not found"
			ings[notFoundCount] = v
			notFoundCount++
		} else if err != nil {
			return nil, err
		}
	}

	if notFoundCount != 0 {
		return ings[:notFoundCount], models.ErrIngsNotFound
	}

	return nil, nil
}

func (core *dalCore) insertToMenuIngs(tx *sqlx.Tx, menuID uint64, ings []models.MenuIngredients) error {
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

	for _, v := range ings {
		v.ProductID = menuID
		_, err = stmt.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
