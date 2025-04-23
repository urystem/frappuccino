package dal

import (
	"time"

	"frappuccino/models"

	"github.com/jmoiron/sqlx"
)

type dalAggregation struct {
	database *sqlx.DB
}

type AggregationDalInter interface {
	AmountSales() (float64, error)
	Popularies() (*models.PopularItems, error)
	CountOfOrderedItems(start, end *time.Time) (map[string]uint64, error)
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

func (db *dalAggregation) Popularies() (*models.PopularItems, error) {
	popularsQ := `
		SELECT oi.product_id, m.name, SUM(oi.quantity) AS sum
			FROM order_items AS oi
			JOIN menu_items AS m ON m.id = oi.product_id
			JOIN orders AS o ON o.id = oi.order_id
			WHERE o.status = 'accepted'
			GROUP BY oi.product_id, m.name
			ORDER BY sum DESC`

	var popularies models.PopularItems
	err := db.database.Select(&popularies.Items, popularsQ)
	return &popularies, err
}

func (db *dalAggregation) CountOfOrderedItems(start, end *time.Time) (map[string]uint64, error) {
	// countItemsQ := `
	// 	SELECT m.name, SUM(oi.quantity) AS sum
	// 		FROM order_items AS oi
	// 		JOIN menu_items AS m ON m.id = oi.product_id
	// 		JOIN orders AS o ON o.id = oi.order_id
	// 		WHERE o.status = 'accepted' AND o.created_at BETWEEN $1 and $2
	// 		GROUP BY m.name
	// 		ORDER BY sum DESC`

	// countItemsQ := `
	// 		SELECT m.name, SUM(oi.quantity) AS sum
	// 			FROM order_items AS oi
	// 			JOIN menu_items AS m ON m.id = oi.product_id
	// 			JOIN orders AS o ON o.id = oi.order_id
	// 			WHERE o.status = 'accepted' AND o.created_at BETWEEN '10-11-2023' and '11-11-2025'
	// 			GROUP BY m.name
	// 			ORDER BY sum DESC`

	countItemsQ2 := `
		SELECT m.name, SUM(oi.quantity) AS sum
			FROM order_items AS oi
			JOIN menu_items AS m ON m.id = oi.product_id
			JOIN orders AS o ON o.id = oi.order_id
			WHERE o.status = 'accepted' AND
				($1::date IS NULL OR o.created_at::date >= $1::date) AND
				($2::date IS NULL OR o.created_at::date <= $2::date)
			GROUP BY m.name
			ORDER BY sum DESC`

	// Было (::date)	Стало (::timestamptz)
	// Усекалась только дата	Учитывается и время
	// 2024-11-10 → 00:00	2024-11-10T13:45:00+06:00
	rows, err := db.database.Query(countItemsQ2, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countItems := make(map[string]uint64) // ✅ предпочтительно
	for rows.Next() {
		var (
			name  string
			count uint64
		)
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		countItems[name] = count
	}
	return countItems, nil
}
