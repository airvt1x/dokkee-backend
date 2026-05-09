//go:build integration

package repository

import (
	"encoding/json"
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResultPostgres_SaveAndGet_Integration(t *testing.T) {
	authRepo := &AuthPostgres{db: testDB}
	user := dokkee.User{
		Username:  "result_user",
		Password:  "hash",
		FirstName: "Res",
		LastName:  "User",
		Email:     "res@test.com",
		Phone:     "+79991112242",
	}
	userID, err := authRepo.CreateUser(user)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM auth_credentials WHERE id = $1", userID) })

	docRepo := &DocumentPostgres{db: testDB}
	doc := dokkee.Document{
		UserID:       userID,
		OriginalName: "result.pdf",
		S3Key:        "result",
		MimeType:     "pdf",
		FileSize:     500,
		Status:       "completed",
	}
	docID, err := docRepo.Create(doc)
	require.NoError(t, err)
	t.Cleanup(func() { testDB.Exec("DELETE FROM documents WHERE id = $1", docID) })

	// Устанавливаем error_msg в пустую строку
	_, err = testDB.Exec("UPDATE documents SET error_msg = '' WHERE id = $1", docID)
	require.NoError(t, err)

	resultRepo := &ResultPostgres{db: testDB}
	resultJSON, _ := json.Marshal(map[string]string{"risk": "low"})
	result := dokkee.AnalysisResult{
		DocumentID: docID,
		ResultJSON: resultJSON,
		ModelUsed:  "test-model",
		TokensUsed: 100,
	}
	err = resultRepo.Save(result)
	assert.NoError(t, err)

	saved, err := resultRepo.GetByDocumentID(docID)
	assert.NoError(t, err)
	assert.Equal(t, docID, saved.DocumentID)
	assert.JSONEq(t, string(resultJSON), string(saved.ResultJSON))
}
