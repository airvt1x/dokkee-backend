package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)


func TestAuthPostgres_UpdateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuthPostgres(sqlxDB)

	userID := 1
	firstName := "Updated"
	lastName := "User"
	middleName := "M"
	phone := "+79991234567"

	input := dokkee.UpdateProfileInput{
		FirstName:  &firstName,
		LastName:   &lastName,
		MiddleName: &middleName,
		Phone:      &phone,
	}

	expectedQuery := regexp.QuoteMeta(`UPDATE user_profiles SET first_name = $1, last_name = $2, middle_name = $3, phone = $4, updated_at = NOW() WHERE user_id = $5`)
	mock.ExpectExec(expectedQuery).
		WithArgs(firstName, lastName, middleName, phone, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateProfile(userID, input)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthPostgres_UpdateProfile_NoFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuthPostgres(sqlxDB)

	input := dokkee.UpdateProfileInput{}

	err = repo.UpdateProfile(1, input)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}