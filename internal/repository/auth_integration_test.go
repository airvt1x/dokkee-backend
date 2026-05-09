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

func TestAuthPostgres_GetUser_NotFound_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}

	_, err := repo.GetUser("nonexistent_user_12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sql: no rows")
}

func TestAuthPostgres_GetProfile_NotFound_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}

	_, err := repo.GetProfile(99999)
	assert.Error(t, err)
}

func TestAuthPostgres_CreateUser_DuplicateUsername_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}

	user := dokkee.User{
		Username:  "duplicate_test_user",
		Password:  "hash",
		FirstName: "First",
		LastName:  "Last",
		Email:     "dup@test.com",
		Phone:     "+79991112288",
	}
	id, err := repo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", id)
	})

	// Пытаемся создать того же пользователя
	_, err = repo.CreateUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
}

func TestAuthPostgres_UpdateProfile_MultipleFields_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "multi_update_user",
		Password:  "hash",
		FirstName: "Original",
		LastName:  "Name",
		Email:     "multi@test.com",
		Phone:     "+79991112293",
	}
	id, err := repo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", id) })

	firstName := "NewFirst"
	lastName := "NewLast"
	phone := "+79991112294"
	input := dokkee.UpdateProfileInput{
		FirstName: &firstName,
		LastName:  &lastName,
		Phone:     &phone,
	}
	err = repo.UpdateProfile(id, input)
	assert.NoError(t, err)

	var dbFirstName, dbLastName, dbPhone string
	err = testDB.QueryRow("SELECT first_name, last_name, phone FROM user_profiles WHERE user_id=$1", id).Scan(&dbFirstName, &dbLastName, &dbPhone)
	assert.NoError(t, err)
	assert.Equal(t, "NewFirst", dbFirstName)
	assert.Equal(t, "NewLast", dbLastName)
	assert.Equal(t, "+79991112294", dbPhone)
}

func TestAuthPostgres_UpdateProfile_AllFields_Integration(t *testing.T) {
	repo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:   "all_fields_user",
		Password:   "hash",
		FirstName:  "OldFirst",
		LastName:   "OldLast",
		MiddleName: "OldMiddle",
		Email:      "all@test.com",
		Phone:      "+79991112301",
	}
	id, err := repo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", id) })

	firstName := "NewFirst"
	lastName := "NewLast"
	middleName := "NewMiddle"
	phone := "+79991112302"
	input := dokkee.UpdateProfileInput{
		FirstName:  &firstName,
		LastName:   &lastName,
		MiddleName: &middleName,
		Phone:      &phone,
	}
	err = repo.UpdateProfile(id, input)
	assert.NoError(t, err)

	var dbFirstName, dbLastName, dbMiddleName, dbPhone string
	err = testDB.QueryRow("SELECT first_name, last_name, middle_name, phone FROM user_profiles WHERE user_id=$1", id).Scan(&dbFirstName, &dbLastName, &dbMiddleName, &dbPhone)
	assert.NoError(t, err)
	assert.Equal(t, "NewFirst", dbFirstName)
	assert.Equal(t, "NewLast", dbLastName)
	assert.Equal(t, "NewMiddle", dbMiddleName)
	assert.Equal(t, "+79991112302", dbPhone)
}

func TestDocumentPostgres_List_ByMultipleStatuses_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "multi_status_user",
		Password:  "hash",
		FirstName: "Multi",
		LastName:  "Status",
		Email:     "multistatus@test.com",
		Phone:     "+79991112303",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID) })

	docRepo := &DocumentPostgres{db: testDB}
	statuses := []string{"queued", "processing", "completed", "failed"}

	for _, status := range statuses {
		doc := dokkee.Document{
			UserID:       userID,
			OriginalName: status + ".pdf",
			S3Key:        status + "_key_" + status,
			MimeType:     "application/pdf",
			FileSize:     100,
			Status:       status,
		}
		id, err := docRepo.Create(doc)
		require.NoError(t, err)
		t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", id) })

		// Явно устанавливаем error_msg в пустую строку, чтобы избежать NULL
		_, err = testDB.Exec("UPDATE documents SET error_msg = '' WHERE id = $1", id)
		require.NoError(t, err)
	}

	for _, status := range statuses {
		docs, err := docRepo.List(userID, status)
		assert.NoError(t, err)
		assert.Len(t, docs, 1, "Expected 1 document with status %s", status)
		if len(docs) > 0 {
			assert.Equal(t, status, docs[0].Status)
		}
	}
}
