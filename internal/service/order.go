package service

import (
	"cafeteria/internal/models"
	"context"
)

// OrderRepository defines database operations related to orders.
type OrderRepository interface {
	GetAll(ctx context.Context) ([]*models.Order, error)
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, order *models.Order) error
	Insert(ctx context.Context, order *models.Order) error
}

// OrderService provides business logic for managing orders.
type OrderService struct {
	Repo OrderRepository
}

// NewOrderService initializes a new order service.
func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{Repo: repo}
}

func (s *OrderService) GetAll(ctx context.Context) ([]*models.Order, error) {
	return s.Repo.GetAll(ctx)
}

func (s *OrderService) GetByID(ctx context.Context, id int) (*models.Order, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *OrderService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

func (s *OrderService) Update(ctx context.Context, order *models.Order) error {
	return s.Repo.Update(ctx, order)
}

func (s *OrderService) Insert(ctx context.Context, order *models.Order) error {
	return s.Repo.Insert(ctx, order)
}
