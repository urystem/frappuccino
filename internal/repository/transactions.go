package repository

import (
	"cafeteria/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type TransactionRepository struct {
	Db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{Db: db}
}

func (r *TransactionRepository) TotalSales(ctx context.Context) (float32, error) {
	var totalSales float32

	err := r.Db.QueryRowContext(ctx, "SELECT COALESCE(SUM(total), 0) FROM orders WHERE status = 'completed'").Scan(&totalSales)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch total sales: %w", err)
	}

	return totalSales, nil
}

func (r *TransactionRepository) PopularItems(ctx context.Context) (models.JSONB, error) {
	var rawJSON []byte // Store raw JSON result

	query := `SELECT jsonb_build_object('popular items', jsonb_agg(jsonb_build_object('name', name, 'total_count', total_count))) AS popular_items
			FROM (
				SELECT 
					mi.name,
					SUM(oi.quantity) AS total_count
				FROM order_items oi
				JOIN menu_items mi ON oi.menu_items_id = mi.menu_items_id
				GROUP BY mi.name
				ORDER BY total_count DESC
			) AS popular_items_subquery;`

	err := r.Db.QueryRowContext(ctx, query).Scan(&rawJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular items: %w", err)
	}

	// Unmarshal raw JSON into a slice of JSONB objects
	var items models.JSONB
	if err := json.Unmarshal(rawJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONB: %w", err)
	}

	return items, nil
}

func (r *TransactionRepository) NumberOfOrderedItems(ctx context.Context, start, end time.Time) (models.JSONB, error) {
	var rawJSON []byte

	var startTime, endTime interface{}

	if start.IsZero() {
		startTime = "2000-01-01"
	} else {
		startTime = start
	}

	if end.IsZero() {
		endTime = time.Now()
	} else {
		endTime = end
	}

	query := `SELECT COALESCE(
					jsonb_object_agg(name, order_count), '{}'
				) AS popular_items
				FROM popular_menu_items
				WHERE last_updated_at BETWEEN $1 AND $2;`

	err := r.Db.QueryRowContext(ctx, query, startTime, endTime).Scan(&rawJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular items: %w", err)
	}

	var items models.JSONB
	if err := json.Unmarshal(rawJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONB: %w", err)
	}

	return items, nil
}

func (r *TransactionRepository) SearchOrders(ctx context.Context, q, filter string, low, high float32) (models.JSONB, error) {
	query := `SELECT * FROM search_all($1, $2, $3, $4);`

	var rawJSON []byte
	err := r.Db.QueryRowContext(ctx, query, q, filter, low, high).Scan(&rawJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular items: %w", err)
	}

	var items models.JSONB
	if err := json.Unmarshal(rawJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONB: %w", err)
	}

	return items, nil
}

func (r *TransactionRepository) OrderedItemsByPeriod(ctx context.Context, period string, month string, year int) (models.JSONB, error) {
	var result models.JSONB
	var rawJSON []byte
	var err error

	switch period {
	case "day":
		query := `
            SELECT COALESCE(
                jsonb_object_agg(day::text, count),
                '{}'::jsonb
            )
            FROM (
                SELECT 
                    EXTRACT(DAY FROM created_at) AS day,
                    COUNT(*) AS count
                FROM (
                    SELECT 
                        o.orders_id,
                        MIN(osh.updated_at) AS created_at
                    FROM orders o
                    JOIN order_status_history osh ON osh.orders_id = o.orders_id
                    WHERE 
                        EXTRACT(MONTH FROM osh.updated_at) = EXTRACT(MONTH FROM TO_DATE($1, 'Month')) AND
                        EXTRACT(YEAR FROM osh.updated_at) = $2
                    GROUP BY o.orders_id
                ) orders_with_dates
                GROUP BY day
                ORDER BY day
            ) AS daily_counts`
		err = r.Db.QueryRowContext(ctx, query, month, year).Scan(&rawJSON)
	case "month":
		query := `
            SELECT COALESCE(
                jsonb_object_agg(
                    LOWER(TO_CHAR(TO_DATE(month::text, 'MM'), 'Month')),
                    count
                ),
                '{}'::jsonb
            )
            FROM (
                SELECT 
                    EXTRACT(MONTH FROM created_at) AS month,
                    COUNT(*) AS count
                FROM (
                    SELECT 
                        o.orders_id,
                        MIN(osh.updated_at) AS created_at
                    FROM orders o
                    JOIN order_status_history osh ON osh.orders_id = o.orders_id
                    WHERE EXTRACT(YEAR FROM osh.updated_at) = $1
                    GROUP BY o.orders_id
                ) orders_with_dates
                GROUP BY month
                ORDER BY month
            ) AS monthly_counts`
		err = r.Db.QueryRowContext(ctx, query, year).Scan(&rawJSON)
	default:
		return nil, fmt.Errorf("invalid period parameter")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch ordered items by period: %w", err)
	}

	if rawJSON == nil {
		return models.JSONB{}, nil
	}

	if err := json.Unmarshal(rawJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ordered items data: %w", err)
	}

	return result, nil
}
