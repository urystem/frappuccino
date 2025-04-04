package dal

import "frappuccino/models"

type OrderDalInter interface {
	SelectAllOrders() ([]models.Order, error)
	SelectOrder(uint64) (*models.Order, error)
	DeleteOrder(uint64) error
}

func (core *dalCore) SelectAllOrders() ([]models.Order, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orders []models.Order

	err = tx.Select(&orders, `SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareNamed(`SELECT product_id, quantity FROM order_items WHERE order_id = :id`)
	if err != nil {
		return nil, nil
	}

	defer stmt.Close()

	for i, v := range orders {
		err = stmt.Select(&orders[i].Items, v)
		if err != nil {
			return nil, err
		}
	}
	return orders, tx.Commit()
}

func (core *dalCore) SelectOrder(id uint64) (*models.Order, error) {
	tx, err := core.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var order models.Order

	err = tx.Get(&order, `SELECT * FROM orders WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	err = tx.Select(&order.Items, `SELECT product_id, quantity FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &order, tx.Commit()
}

func (core *dalCore) DeleteOrder(id uint64) error {
	tx, err := core.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.Exec(`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return err
	}

	affects, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affects == 0 {
		return models.ErrNotFound
	}

	return tx.Commit()
}
