package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"  

	"github.com/DATA-DOG/go-sqlmock"
	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestResultPostgres_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewResultPostgres(sqlxDB)

	resultJSON, _ := json.Marshal(map[string]string{"risk": "high"})
	result := dokkee.AnalysisResult{
		DocumentID: 10,
		ResultJSON: resultJSON,
		ModelUsed:  "deepseek-chat",
		TokensUsed: 150,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO analysis_results (document_id, result_json, model_used, tokens_used) VALUES ($1, $2, $3, $4)`)).
		WithArgs(result.DocumentID, result.ResultJSON, result.ModelUsed, result.TokensUsed).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(result)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResultPostgres_GetByDocumentID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewResultPostgres(sqlxDB)

	docID := 10
	now := time.Now().Round(time.Second)
	resultJSON := json.RawMessage(`{"risk":"high"}`)

	rows := sqlmock.NewRows([]string{"id", "document_id", "result_json", "model_used", "tokens_used", "analyzed_at"}).
		AddRow(1, 10, resultJSON, "deepseek-chat", 150, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM analysis_results WHERE document_id = $1`)).
		WithArgs(docID).
		WillReturnRows(rows)

	result, err := repo.GetByDocumentID(docID)
	assert.NoError(t, err)
	assert.Equal(t, 10, result.DocumentID)
	assert.Equal(t, resultJSON, result.ResultJSON)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResultPostgres_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewResultPostgres(sqlxDB)

	resultJSON, _ := json.Marshal(map[string]string{"risk": "high"})
	result := dokkee.AnalysisResult{
		DocumentID: 10,
		ResultJSON: resultJSON,
		ModelUsed:  "deepseek-chat",
		TokensUsed: 150,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO analysis_results (document_id, result_json, model_used, tokens_used) VALUES ($1, $2, $3, $4)`)).
		WithArgs(result.DocumentID, result.ResultJSON, result.ModelUsed, result.TokensUsed).
		WillReturnError(errors.New("duplicate key"))

	err = repo.Save(result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResultPostgres_GetByDocumentID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewResultPostgres(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM analysis_results WHERE document_id = $1`)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByDocumentID(999)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}