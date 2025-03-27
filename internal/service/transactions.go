package service

import (
	"context"
)

type TransactionRepository interface {
	TotalSales(ctx context.Context) (float32, error)
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
