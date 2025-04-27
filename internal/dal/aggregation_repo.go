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
	SearchByWordInventory(ind string, minPrice, maxPrice float64, stc *models.SearchThings) error
	SearchByWordMenu(find string, minPrice, maxPrice float64, strc *models.SearchThings) error
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

func (db *dalAggregation) SearchByWordInventory(find string, minPrice, maxPrice float64, strc *models.SearchThings) error {
	query := `
	WITH ranked_inventory AS (
    	SELECT
        	id, name, description, quantity,
        	reorder_level, unit, price,
        	ROUND(ts_rank(
            	setweight(to_tsvector(name),'A') ||
            	setweight(to_tsvector(description), 'B'),
            	to_tsquery($1)
        	)::numeric, 2) AS relevance
    	FROM inventory
		WHERE price BETWEEN $2 AND $3
	)
	SELECT *
	FROM ranked_inventory
	WHERE relevance > 0.009
	ORDER BY relevance DESC`

	return db.database.Select(&strc.Inventories, query, find, minPrice, maxPrice)
}

func (db *dalAggregation) SearchByWordMenu(find string, minPrice, maxPrice float64, strc *models.SearchThings) error {
	query := `
	WITH ranked_menu AS(
		SELECT
			m.id,
			m.name,
			m.description,
			m.tags,
			m.allergens,
			m.price,
			array_agg(i.name) AS inventories,
			ROUND(
				ts_rank(
					setweight(to_tsvector(m.name),'A') ||
					setweight(to_tsvector(m.description),'B') ||
					setweight(to_tsvector(array_to_string(m.tags, ' ')), 'C') ||
					setweight(to_tsvector(array_to_string(m.allergens, ' ')), 'D') ||
					setweight(to_tsvector(string_agg(i.name, ' ')), 'B'),
					to_tsquery($1)
				)::numeric, 2) AS relevance
		FROM menu_items AS m
		JOIN menu_item_ingredients AS mi ON m.id=mi.product_id
		JOIN inventory AS i ON mi.inventory_id = i.id
		WHERE m.price BETWEEN $2 AND $3
		GROUP BY m.id
	)
	SELECT * 
	FROM ranked_menu 
	WHERE relevance > 0.009
	ORDER BY relevance DESC`
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = tx.Select(&strc.Menus, query, find, minPrice, maxPrice)
	if err != nil {
		return err
	}

	// stmt, err := tx.PrepareNamed(`SELECT inventory_id, quantity FROM menu_item_ingredients WHERE product_id=:id`)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()

	// for i, menu := range strc.Menus {
	// 	err = stmt.Select(&strc.Menus[i].Ingredients, menu)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return tx.Commit()
}
