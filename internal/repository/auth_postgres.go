package repository

import (
	"fmt"

	"github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user dokkee.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, first_name, last_name, middle_name, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Password, user.FirstName, user.LastName, user.MiddleName, user.Email, user.Phone)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (dokkee.User, error) {
	var user dokkee.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}
