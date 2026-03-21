package repository

import (
	"github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user dokkee.User) (int, error)
	GetUser(email, password string) (dokkee.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
