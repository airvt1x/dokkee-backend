package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuthPostgres(sqlxDB)

	user := dokkee.User{
		Username:   "testuser",
		Password:   "hashedpassword",
		FirstName:  "Test",
		LastName:   "User",
		MiddleName: "Middle",
		Email:      "test@example.com",
		Phone:      "+1234567890",
	}

	mock.ExpectQuery(`INSERT INTO users \(username, password_hash, first_name, last_name, middle_name, email, phone\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id`).
		WithArgs(user.Username, user.Password, user.FirstName, user.LastName, user.MiddleName, user.Email, user.Phone).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuthPostgres(sqlxDB)

	email := "test@example.com"
	password := "hashedpassword"

	expectedUser := dokkee.User{
		Id: 1,
	}

	mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1 AND password_hash = \$2`).
		WithArgs(email, password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user, err := repo.GetUser(email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}
