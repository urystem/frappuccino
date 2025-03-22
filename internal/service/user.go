package service

import (
	"cafeteria/internal/models"
	"context"
)

type UserRepository interface {
	Register(ctx context.Context, user *models.User) (string, error)
	GetToken(ctx context.Context, username, pass string) (string, error)
}

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(ctx context.Context, user *models.User) (string, error) {
	return s.Repo.Register(ctx, user)
}

func (s *UserService) GetToken(ctx context.Context, username, pass string) (string, error) {
	return s.Repo.GetToken(ctx, username, pass)
}
