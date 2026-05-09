//go:build integration

package repository

import (
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthPostgres_CreateUser_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:   "int_test_user",
		Password:   "hashedpass",
		FirstName:  "Integration",
		LastName:   "Test",
		MiddleName: "User",
		Email:      "int@test.com",
		Phone:      "+79991112233",
	}
	id, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", id)
	})
}

func TestAuthPostgres_GetUser_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "getuser_int",
		Password:  "hash",
		FirstName: "Get",
		LastName:  "User",
		Email:     "get@test.com",
		Phone:     "+79991112234",
	}
	id, err := repo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", id)
	})

	found, err := repo.GetUser("getuser_int")
	assert.NoError(t, err)
	assert.Equal(t, id, found.Id)
	// Username не маппится в структуру, поэтому не проверяем
}

func TestAuthPostgres_UpdateProfile_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "update_int",
		Password:  "hash",
		FirstName: "Old",
		LastName:  "Name",
		Email:     "update@test.com",
		Phone:     "+79991112236",
	}
	id, err := repo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", id)
	})

	newFirstName := "NewFirstName"
	input := dokkee.UpdateProfileInput{FirstName: &newFirstName}
	err = repo.UpdateProfile(id, input)
	assert.NoError(t, err)

	var firstName string
	err = testDB.Get(&firstName, "SELECT first_name FROM user_profiles WHERE user_id=$1", id)
	assert.NoError(t, err)
	assert.Equal(t, "NewFirstName", firstName)
}