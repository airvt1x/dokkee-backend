package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dokkee "github.com/airvt1x/dokkee-backend"
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

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO auth_credentials`).
		WithArgs(user.Username, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO user_profiles`).
		WithArgs(1, user.FirstName, user.LastName, user.MiddleName, user.Email, user.Phone).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO user_balances`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

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

	username := "testuser"

	mock.ExpectQuery(`SELECT id, password_hash AS password FROM auth_credentials WHERE username = \$1`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(1, "hashedpassword"))

	user, err := repo.GetUser(username)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}
