package service

import (
	"cafeteria/internal/models"
	"context"
	"errors"
)

type MenuRepository interface {
	GetAll(ctx context.Context) ([]models.Menu, error)
	GetElementById(ctx context.Context, MenuId int) (models.Menu, error)
	Delete(ctx context.Context, MenuId int) error
	Put(ctx context.Context, item models.Menu) error
	Post(ctx context.Context, item models.Menu) error
}

type MenuService struct {
	Repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{Repo: repo}
}

func (s *MenuService) GetAll(ctx context.Context) ([]models.Menu, error) {
	items, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *MenuService) GetElementById(ctx context.Context, MenuId int) (models.Menu, error) {
	item, err := s.Repo.GetElementById(ctx, MenuId)
	if err != nil {
		return models.Menu{}, err
	}
	return item, nil
}

func (s *MenuService) Put(ctx context.Context, item models.Menu) error {
	if item.MenuID <= 0 {
		err := errors.New("invalid item ID")
		return err
	}

	err := s.Repo.Put(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (s *MenuService) Delete(ctx context.Context, MenuId int) error {
	err := s.Repo.Delete(ctx, MenuId)
	if err != nil {
		return err
	}
	return nil
}

func (s *MenuService) Post(ctx context.Context, item models.Menu) error {
	if item.Name == "" || item.Price < 0 || item.Description == "" {
		err := errors.New("invalid item data")
		return err
	}

	err := s.Repo.Post(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
