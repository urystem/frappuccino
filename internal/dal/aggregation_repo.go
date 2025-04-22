package dal

import (
	"frappuccino/models"

	"github.com/jmoiron/sqlx"
)

type dalAggregation struct {
	database *sqlx.DB
}

type AggregationDalInter interface {
	AmountSales() (float64, error)
	Popularies() ([]models.PopularItems, error)
}

func ReturnDulAggregationDB(db *sqlx.DB) AggregationDalInter {
	return &dalAggregation{db}
}

func (db *dalAggregation) AmountSales() (float64, error) {
	sumTotal := `
	SELECT SUM(total)
	FROM orders
	WHERE status = 'accepted'`

	var total float64
	tx, err := db.database.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	err = tx.Get(&total, sumTotal)
	if err != nil {
		return 0, err
	}
	return total, tx.Commit()
}

func (db *dalAggregation) Popularies() ([]models.PopularItems, error) {
	popularsQ := `
		SELECT oi.product_id, m.name, SUM(oi.quantity) AS sum
			FROM order_items AS oi
			JOIN menu_items AS m ON m.id = oi.product_id
			JOIN orders AS o ON o.id = oi.order_id
			WHERE o.status = 'accepted'
			GROUP BY oi.product_id, m.name
			ORDER BY sum DESC`

	var popularies []models.PopularItems
	err := db.database.Select(&popularies, popularsQ)
	return popularies, err
}

