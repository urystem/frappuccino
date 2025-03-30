package service

import (
	"cafeteria/internal/models"
	"context"
	"time"
)

type TransactionRepository interface {
	TotalSales(ctx context.Context) (float32, error)
	PopularItems(ctx context.Context) (models.JSONB, error)
	NumberOfOrderedItems(ctx context.Context, start, end time.Time) (models.JSONB, error)
	SearchOrders(ctx context.Context, q, filter string, low, high float32) (models.JSONB, error)
	OrderedItemsByPeriod(ctx context.Context, period string, month string, year int) (models.JSONB, error)
}

// TransactionService provides business logic for managing orders.
type TransactionService struct {
	Repo TransactionRepository
}

// NewTransactionService initializes a new order service.
func NewTransactionService(repo TransactionRepository) *TransactionService {
	return &TransactionService{Repo: repo}
}

func (s *TransactionService) TotalSales(ctx context.Context) (float32, error) {
	total, err := s.Repo.TotalSales(ctx)
	if err != nil {
		return 0.0, err
	}
	return total, nil
}

func (s *TransactionService) PopularItems(ctx context.Context) (models.JSONB, error) {
	return s.Repo.PopularItems(ctx)
}

func (s *TransactionService) NumberOfOrderedItems(ctx context.Context, start, end time.Time) (models.JSONB, error) {
	return s.Repo.NumberOfOrderedItems(ctx, start, end)
}

func (s *TransactionService) SearchOrders(ctx context.Context, q, filter string, low, high float32) (models.JSONB, error) {
	return s.Repo.SearchOrders(ctx, q, filter, low, high)
}

func (s *TransactionService) OrderedItemsByPeriod(ctx context.Context, period string, month string, year int) (models.JSONB, error) {
	return s.Repo.OrderedItemsByPeriod(ctx, period, month, year)
}
