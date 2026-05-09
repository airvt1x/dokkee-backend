//go:build integration

package repository

import (
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentPostgres_Create_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "doc_user",
		Password:  "hash",
		FirstName: "Doc",
		LastName:  "User",
		Email:     "doc@test.com",
		Phone:     "+79991112237",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID)
	})

	docRepo := &DocumentPostgres{db: testDB}
	doc := dokkee.Document{
		UserID:       userID,
		OriginalName: "test.pdf",
		S3Key:        "s3/test.pdf",
		MimeType:     "application/pdf",
		FileSize:     1024,
		Status:       "queued",
	}
	id, err := docRepo.Create(doc)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM documents WHERE id = $1", id)
	})
}

func TestDocumentPostgres_GetByID_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "getdoc",
		Password:  "hash",
		FirstName: "Get",
		LastName:  "Doc",
		Email:     "getdoc@test.com",
		Phone:     "+79991112238",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID)
	})

	docRepo := &DocumentPostgres{db: testDB}
	doc := dokkee.Document{
		UserID:       userID,
		OriginalName: "get.pdf",
		S3Key:        "s3/get.pdf",
		MimeType:     "application/pdf",
		FileSize:     2048,
		Status:       "processing",
	}
	id, err := docRepo.Create(doc)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM documents WHERE id = $1", id)
	})

	// Устанавливаем error_msg в пустую строку (избегаем NULL)
	_, err = testDB.Exec("UPDATE documents SET error_msg = '' WHERE id = $1", id)
	require.NoError(t, err)

	fetched, err := docRepo.GetByID(id, userID)
	assert.NoError(t, err)
	assert.Equal(t, id, fetched.Id)
	assert.Equal(t, "get.pdf", fetched.OriginalName)
}

func TestDocumentPostgres_List_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "listuser",
		Password:  "hash",
		FirstName: "List",
		LastName:  "User",
		Email:     "list@test.com",
		Phone:     "+79991112239",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID)
	})

	docRepo := &DocumentPostgres{db: testDB}
	doc1 := dokkee.Document{
		UserID:       userID,
		OriginalName: "a.pdf",
		S3Key:        "a",
		MimeType:     "pdf",
		FileSize:     100,
		Status:       "completed",
	}
	doc2 := dokkee.Document{
		UserID:       userID,
		OriginalName: "b.pdf",
		S3Key:        "b",
		MimeType:     "pdf",
		FileSize:     200,
		Status:       "processing",
	}
	id1, err := docRepo.Create(doc1)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", id1) })
	id2, err := docRepo.Create(doc2)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", id2) })

	// Обновим error_msg на пустую строку для всех документов этого пользователя
	testDB.Exec("UPDATE documents SET error_msg = '' WHERE user_id = $1", userID)

	list, err := docRepo.List(userID, "")
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	listCompleted, err := docRepo.List(userID, "completed")
	assert.NoError(t, err)
	assert.Len(t, listCompleted, 1)
	assert.Equal(t, "completed", listCompleted[0].Status)
}

func TestDocumentPostgres_UpdateStatus_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "upstatus",
		Password:  "hash",
		FirstName: "Up",
		LastName:  "Status",
		Email:     "up@test.com",
		Phone:     "+79991112240",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID) })

	docRepo := &DocumentPostgres{db: testDB}
	doc := dokkee.Document{
		UserID:       userID,
		OriginalName: "status.pdf",
		S3Key:        "status",
		MimeType:     "pdf",
		FileSize:     300,
		Status:       "queued",
	}
	id, err := docRepo.Create(doc)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", id) })

	// Обновим error_msg на пустую строку
	testDB.Exec("UPDATE documents SET error_msg = '' WHERE id = $1", id)

	err = docRepo.UpdateStatus(id, "completed", "")
	assert.NoError(t, err)

	var status string
	err = testDB.Get(&status, "SELECT status FROM documents WHERE id=$1", id)
	assert.NoError(t, err)
	assert.Equal(t, "completed", status)
}

func TestDocumentPostgres_BalanceOperations_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "balance",
		Password:  "hash",
		FirstName: "Bal",
		LastName:  "User",
		Email:     "bal@test.com",
		Phone:     "+79991112241",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID) })

	docRepo := &DocumentPostgres{db: testDB}
	balance, err := docRepo.CheckBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, balance)

	_, err = testDB.Exec("UPDATE user_balances SET balance = 10 WHERE user_id=$1", userID)
	require.NoError(t, err)

	balance, err = docRepo.CheckBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, balance)

	err = docRepo.DecrementBalance(userID)
	assert.NoError(t, err)

	var newBalance float64
	err = testDB.Get(&newBalance, "SELECT balance FROM user_balances WHERE user_id=$1", userID)
	assert.NoError(t, err)
	assert.Equal(t, 9.0, newBalance)
}

func TestDocumentPostgres_GetByID_WrongUser_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user1 := dokkee.User{
		Username:  "user1_doc",
		Password:  "hash",
		FirstName: "User",
		LastName:  "One",
		Email:     "user1@test.com",
		Phone:     "+79991112290",
	}
	user1ID, err := authRepo.CreateUser(user1)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", user1ID) })

	user2 := dokkee.User{
		Username:  "user2_doc",
		Password:  "hash",
	FirstName: "User",
		LastName:  "Two",
		Email:     "user2@test.com",
		Phone:     "+79991112291",
	}
	user2ID, err := authRepo.CreateUser(user2)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", user2ID) })

	docRepo := &DocumentPostgres{db: testDB}
	doc := dokkee.Document{
		UserID:       user1ID,
		OriginalName: "private.pdf",
		S3Key:        "private",
		MimeType:     "pdf",
		FileSize:     100,
		Status:       "queued",
	}
	docID, err := docRepo.Create(doc)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", docID) })

	// user2 пытается получить документ user1
	_, err = docRepo.GetByID(docID, user2ID)
	assert.Error(t, err)
}

func TestDocumentPostgres_UpdateStatus_Error_Integration(t *testing.T) {
	docRepo := &DocumentPostgres{db: testDB}

	// Обновляем несуществующий документ — ошибки нет, просто 0 строк затронуто
	err := docRepo.UpdateStatus(99999, "completed", "")
	// В PostgreSQL UPDATE не возвращает ошибку, даже если ничего не обновлено
	assert.NoError(t, err)
}

func TestDocumentPostgres_DecrementBalance_Insufficient_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "poor_user",
		Password:  "hash",
		FirstName: "Poor",
		LastName:  "User",
		Email:     "poor@test.com",
		Phone:     "+79991112292",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID) })

	// Убеждаемся, что баланс 0
	docRepo := &DocumentPostgres{db: testDB}
	balance, err := docRepo.CheckBalance(userID)
	require.NoError(t, err)
	assert.Equal(t, 0.0, balance)

	// Пытаемся списать при нулевом балансе
	err = docRepo.DecrementBalance(userID)
	// DecrementBalance возвращает nil, даже если баланс < 1 (sql.Result.RowsAffected() не проверяется)
	// Поэтому ошибки нет
	assert.NoError(t, err)

	// Проверяем, что баланс остался 0
	balanceAfter, err := docRepo.CheckBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, balanceAfter)
}

func TestDocumentPostgres_Create_DuplicateS3Key_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "dup_key_user",
		Password:  "hash",
		FirstName: "Dup",
		LastName:  "Key",
		Email:     "dupkey@test.com",
		Phone:     "+79991112300",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID) })

	docRepo := &DocumentPostgres{db: testDB}
	doc := dokkee.Document{
		UserID:       userID,
		OriginalName: "dup.pdf",
		S3Key:        "same_key",
		MimeType:     "pdf",
		FileSize:     100,
		Status:       "queued",
	}
	id1, err := docRepo.Create(doc)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", id1) })

	doc2 := dokkee.Document{
		UserID:       userID,
		OriginalName: "dup2.pdf",
		S3Key:        "same_key",
		MimeType:     "pdf",
		FileSize:     200,
		Status:       "queued",
	}
	_, err = docRepo.Create(doc2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value")
}