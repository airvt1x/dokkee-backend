package service

import (
	"github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/repository"
)

type Authorization interface {
	CreateUser(user dokkee.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
