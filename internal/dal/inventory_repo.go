package dal

import "hot-coffee/models"

type InventoryDataAccess interface {
	// InsertInventoryV1(*models.Inventory) error
	// InsertInventoryV2(*models.Inventory) error
	// InsertInventoryV3(*models.Inventory) error
	// InsertInventoryV4(*models.Inventory) error
	InsertInventoryV5(*models.Inventory) error
	// InsertInventoryV6(*models.Inventory) error
	SelectInventories() ([]models.Inventory, error)
	SelectInventory(uint64) (*models.Inventory, error)
	UpdateInventory(*models.Inventory) error
	DeleteInventory(uint64) (*models.InventoryDepend, error)
}

func (core *dalCore) InsertInventoryV1(inv *models.Inventory) error {
	_, err := core.db.Exec(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES($1,$2,$3,$4,$5,$6)`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price,
	)
	return err
}

func (core *dalCore) InsertInventoryV2(inv *models.Inventory) error {
	_, err := core.db.NamedExec(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES (:name, :description, :quantity, :reorder_level, :unit, :price)
	`, inv)
	return err
}

func (core *dalCore) InsertInventoryV3(inv *models.Inventory) error {
	return core.db.QueryRow(`
	INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
		VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price).Scan(&inv.ID)
}

func (core *dalCore) InsertInventoryV4(inv *models.Inventory) error {
	return core.db.Get(&inv.ID, `
    INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id`, inv.Name, inv.Descrip, inv.Quantity, inv.ReorderLvl, inv.Unit, inv.Price)
}

func (core *dalCore) InsertInventoryV5(inv *models.Inventory) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// tx.QueryRowx также подходит
	if err = tx.QueryRow(`
		INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
			VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id`,
		inv.Name,
		inv.Descrip,
		inv.Quantity,
		inv.ReorderLvl,
		inv.Unit,
		inv.Price).Scan(&inv.ID); err != nil {
		return err
	}

	if _, err = tx.NamedExec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES (:id, :quantity, 'restock')
	`, inv); err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalCore) InsertInventoryV6(inv *models.Inventory) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareNamed(`
	INSERT INTO inventory (name, description, quantity, reorder_level, unit, price)
	VALUES (:name, :description, :quantity, :reorder_level, :unit, :price)
	RETURNING id`)
	if err != nil {
		return err // Ошибка при подготовке запроса
	}
	defer stmt.Close()

	// Выполнение подготовленного запроса
	if err = stmt.QueryRow(inv).Scan(&inv.ID); err != nil {
		return err // Ошибка при выполнении
	}

	if _, err = tx.NamedExec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES (:id, :quantity, 'restock')
	`, inv); err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalCore) SelectInventories() ([]models.Inventory, error) {
	var invts []models.Inventory
	err := core.db.Select(&invts, "SELECT * FROM inventory")
	return invts, err
}

func (core *dalCore) SelectInventory(id uint64) (*models.Inventory, error) {
	var inv models.Inventory
	return &inv, core.db.Get(&inv, "SELECT * FROM inventory WHERE id = $1", id)
}

func (core *dalCore) UpdateInventory(inv *models.Inventory) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var quantity_changed float64
	// err = tx.QueryRow(`SELECT quantity FROM inventory WHERE id=$1`,inv.ID).Scan(&oldQuantity)
	if err = tx.Get(&quantity_changed, `SELECT quantity FROM inventory WHERE id=$1`, inv.ID); err != nil {
		return err
	}

	res, err := tx.NamedExec(`
	UPDATE inventory
		SET name = :name, description = :description, quantity = :quantity,
		    reorder_level = :reorder_level, unit = :unit, price = :price
		WHERE id = :id`, inv)
	if err != nil {
		return err
	}

	// UPDATE егер id табылмаса да қате қайтармайды.
	// сондықтан кестенің өзгерісін rowsAffected пен тексереміз
	// егер айди бар болып, бірақ жаңа кестеден айырмасы жоқ болса да, rowsAffected =1 болады
	// демек тек жоқ айдиде ғана 0 болады
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	quantity_changed = inv.Quantity - quantity_changed

	reason := "restock"
	if quantity_changed < 0 {
		reason = "annul"
	}

	if _, err = tx.Exec(`
	INSERT INTO inventory_transactions (inventory_id, quantity_change, reason)
		VALUES ($1, $2, $3)
	`, inv.ID, quantity_changed, reason); err != nil {
		return err
	}
	return tx.Commit()
}

func (core *dalCore) DeleteInventory(id uint64) (*models.InventoryDepend, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var menuDepend models.InventoryDepend

	menusNames := `SELECT id, name
		FROM menu_items
		JOIN menu_item_ingredients ON id=product_id
		WHERE inventory_id=$1`

	err = tx.Select(&menuDepend.Menus, menusNames, id)
	if err != nil {
		return nil, err
	}

	if len(menuDepend.Menus) != 0 {
		menuDepend.Err = models.ErrDelDepend.Error()
		return &menuDepend, nil
	}
	res, err := tx.Exec(`DELETE FROM inventory WHERE id = $1`, id)
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
