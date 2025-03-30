package service

import (
	"cafeteria/internal/models"
	"context"
)

// InventoryRepository defines database operations related to inventory.
type InventoryRepository interface {
	GetAll(ctx context.Context) ([]*models.InventoryItem, error)
	GetByID(ctx context.Context, id int) (*models.InventoryItem, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, item *models.InventoryItem) error
	Insert(ctx context.Context, item *models.InventoryItem) error
	GetLeftovers(ctx context.Context, sortBy string, page, pageSize int) (models.LeftoversResponse, error)
}

// InventoryService provides business logic for managing inventory.
type InventoryService struct {
	Repo InventoryRepository
}

// NewInventoryService initializes a new inventory service.
func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{Repo: repo}
}

func (s *InventoryService) GetAll(ctx context.Context) ([]*models.InventoryItem, error) {
	return s.Repo.GetAll(ctx)
}

func (s *InventoryService) GetByID(ctx context.Context, id int) (*models.InventoryItem, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *InventoryService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

func (s *InventoryService) Update(ctx context.Context, item *models.InventoryItem) error {
	if err := item.IsValid(); err != nil {
		return err
	}
	return s.Repo.Update(ctx, item)
}

func (s *InventoryService) Insert(ctx context.Context, item *models.InventoryItem) error {
	if err := item.IsValid(); err != nil {
		return err
	}
	return s.Repo.Insert(ctx, item)
}

func (s *InventoryService) GetLeftovers(ctx context.Context, sortBy string, page, pageSize int) (models.LeftoversResponse, error) {
	return s.Repo.GetLeftovers(ctx, sortBy, page, pageSize)
}
