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

	query := `SELECT COALESCE(
					jsonb_object_agg(name, order_count), '{}'
				) AS popular_items
				FROM popular_menu_items`

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

	fmt.Println(start)
	fmt.Println(end)

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
