package repository

import (
	"context"
	"database/sql"
	"fmt"
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
