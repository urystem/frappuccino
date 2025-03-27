package service

import (
	"cafeteria/internal/models"
	"context"
)

// MenuRepository defines the methods required for menu data operations.
type MenuRepository interface {
	GetAll(ctx context.Context) ([]*models.MenuItem, error)
	GetByID(ctx context.Context, id int) (*models.MenuItem, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, item *models.MenuItem) error
	Insert(ctx context.Context, item *models.MenuItem) error
}

type MenuService struct {
	Repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{Repo: repo}
}

func (s *MenuService) GetAll(ctx context.Context) ([]*models.MenuItem, error) {
	return s.Repo.GetAll(ctx)
}

func (s *MenuService) GetByID(ctx context.Context, id int) (*models.MenuItem, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *MenuService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

func (s *MenuService) Update(ctx context.Context, item *models.MenuItem) error {
	if err := item.IsValid(); err != nil {
		return err
	}
	return s.Repo.Update(ctx, item)
}

func (s *MenuService) Insert(ctx context.Context, item *models.MenuItem) error {
	if err := item.IsValid(); err != nil {
		return err
	}
	return s.Repo.Insert(ctx, item)
}
