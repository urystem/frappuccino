package service

import (
	"cafeteria/internal/models"
	"context"
)

type InventoryRepository interface {
	GetAll(ctx context.Context) ([]*models.InventoryItem, error)
	GetByID(ctx context.Context, id int) (*models.InventoryItem, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, item *models.InventoryItem) error
	Insert(ctx context.Context, item *models.InventoryItem) error
}

type InventoryService struct {
	Repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{Repo: repo}
}

func (s *InventoryService) GetAll(ctx context.Context) ([]*models.InventoryItem, error) {
	inventory, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *InventoryService) GetByID(ctx context.Context, id int) (*models.InventoryItem, error) {
	inventory, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *InventoryService) Delete(ctx context.Context, id int) error {
	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *InventoryService) Update(ctx context.Context, item *models.InventoryItem) error {
	if err := s.Repo.Update(ctx, item); err != nil {
		return err
	}
	return nil
}

func (s *InventoryService) Insert(ctx context.Context, item *models.InventoryItem) error {
	if err := s.Repo.Insert(ctx, item); err != nil {
		return err
	}
	return nil
}
