package service

import (
	"cafeteria/internal/models"
	"context"
	"errors"
)

type InventoryRepository interface {
	GetAll(ctx context.Context) ([]models.Inventory, error)
	GetElementById(ctx context.Context, InventoryId int) (models.Inventory, error)
	Delete(ctx context.Context, InventoryId int) error
	Put(ctx context.Context, item models.Inventory) error
	Post(ctx context.Context, item models.Inventory) error
}

type InventoryService struct {
	Repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{Repo: repo}
}

func (s *InventoryService) GetAll(ctx context.Context) ([]models.Inventory, error) {
	items, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *InventoryService) GetElementById(ctx context.Context, InventoryId int) (models.Inventory, error) {
	item, err := s.Repo.GetElementById(ctx, InventoryId)
	if err != nil {
		return models.Inventory{}, err
	}
	return item, nil
}

func (s *InventoryService) Put(ctx context.Context, item models.Inventory) error {
	if item.InventoryId <= 0 {
		err := errors.New("invalid item ID")
		return err
	}

	err := s.Repo.Put(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (s *InventoryService) Delete(ctx context.Context, InventoryId int) error {
	err := s.Repo.Delete(ctx, InventoryId)
	if err != nil {
		return err
	}
	return nil
}

func (s *InventoryService) Post(ctx context.Context, item models.Inventory) error {
	if item.Name == "" || item.Quantity < 0 {
		err := errors.New("invalid item data")
		return err
	}

	err := s.Repo.Post(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
