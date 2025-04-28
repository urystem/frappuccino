package dal

import (
	"database/sql"
	"fmt"
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
	SearchByWordOrder(find string, minPrice, maxPrice float64, strc *models.SearchThings) error
	PeriodMonth(month time.Month) ([]map[string]uint64, error)
	PeriodYear(int) ([]map[string]uint64, error)
	GetLeftOversRepo(*models.GetLeftOvers) error
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
	return total, db.database.Get(&total, sumTotal)
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

	return &popularies, db.database.Select(&popularies.Items, popularsQ)
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
	return db.database.Select(&strc.Menus, query, find, minPrice, maxPrice)
}

func (db *dalAggregation) SearchByWordOrder(find string, minPrice, maxPrice float64, strc *models.SearchThings) error {
	query := `
	WITH ranked_order AS(
		SELECT 
			o.id,
			o.customer_name,
			o.status,
			o.allergens,
			o.total,
			array_agg(m.name) AS menu_items,
			ROUND(
				ts_rank(
					setweight(to_tsvector(o.customer_name),'A') ||
					setweight(to_tsvector(array_to_string(o.allergens, ' ')), 'C') ||
					setweight(to_tsvector(string_agg(m.name, ' ')), 'B'),
					to_tsquery($1)
				)::numeric, 2) AS relevance
		FROM orders AS o
		JOIN order_items AS oi ON o.id = oi.order_id
		JOIN menu_items AS m ON oi.product_id = m.id
		WHERE o.total BETWEEN $2 AND $3
		GROUP BY o.id
	)
	SELECT * 
	FROM ranked_order
	WHERE relevance > 0.009
	ORDER BY relevance DESC`

	return db.database.Select(&strc.Orders, query, find, minPrice, maxPrice)
}

func (db *dalAggregation) PeriodMonth(month time.Month) ([]map[string]uint64, error) {
	query := `
		SELECT
			EXTRACT(DAY FROM created_at) AS day,
	        COUNT(*) AS total_orders
		FROM
			orders
		WHERE
			EXTRACT(MONTH FROM created_at) = $1
		GROUP BY day
		ORDER BY day`
	rows, err := db.database.Query(query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return db.rowsToMap(rows)
}

func (db *dalAggregation) PeriodYear(year int) ([]map[string]uint64, error) {
	query := `
	SELECT 
		TO_CHAR(created_at, 'FMMonth') AS month,
		COUNT(*) AS total_orders
	FROM orders
	WHERE 
		EXTRACT(YEAR FROM created_at) = $1
		AND status = 'accepted'
	GROUP BY month, EXTRACT(MONTH FROM created_at)
	ORDER BY EXTRACT(MONTH FROM created_at)`
	rows, err := db.database.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return db.rowsToMap(rows)
}

func (db *dalAggregation) rowsToMap(rows *sql.Rows) ([]map[string]uint64, error) {
	var result []map[string]uint64
	for rows.Next() {
		var day string
		var totalOrders uint64
		err := rows.Scan(&day, &totalOrders)
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]uint64{day: totalOrders})
	}
	return result, nil
}

func (db *dalAggregation) GetLeftOversRepo(over *models.GetLeftOvers) error {
	tx, err := db.database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var countRows uint64
	err = tx.Get(&countRows, `SELECT COUNT(*) FROM inventory`)
	if err != nil {
		return err
	}
	fmt.Println(countRows)
	query := `
		SELECT 
			id, name, quantity, price
		FROM inventory
		ORDER BY  
			CASE 
				WHEN LOWER($1) = 'quantity' THEN quantity
				WHEN LOWER($1) = 'price' THEN price
			 	ELSE id
			END 
			ASC
		LIMIT $2 OFFSET $3`
	offset := (over.CurrentPage - 1) * over.PageSize
	err = tx.Select(&over.Data, query, over.SortBy, over.PageSize, offset)
	if err != nil {
		return err
	}
	over.TotalPages = countRows / over.PageSize
	return tx.Commit()
}
