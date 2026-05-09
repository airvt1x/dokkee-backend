package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestDocumentPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	doc := dokkee.Document{
		UserID:       1,
		OriginalName: "test.pdf",
		S3Key:        "documents/uuid/test.pdf",
		MimeType:     "application/pdf",
		FileSize:     1024,
		Status:       "queued",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO documents (user_id, original_name, s3_key, mime_type, file_size, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)).
		WithArgs(doc.UserID, doc.OriginalName, doc.S3Key, doc.MimeType, doc.FileSize, doc.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	id, err := repo.Create(doc)
	assert.NoError(t, err)
	assert.Equal(t, 10, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentPostgres_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	docID := 5
	userID := 1
	now := time.Now().Round(time.Second)

	rows := sqlmock.NewRows([]string{"id", "user_id", "original_name", "s3_key", "mime_type", "file_size", "status", "error_msg", "created_at", "updated_at"}).
		AddRow(5, 1, "doc.pdf", "s3key", "application/pdf", 2048, "completed", "", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM documents WHERE id = $1 AND user_id = $2`)).
		WithArgs(docID, userID).
		WillReturnRows(rows)

	doc, err := repo.GetByID(docID, userID)
	assert.NoError(t, err)
	assert.Equal(t, 5, doc.Id)
	assert.Equal(t, "doc.pdf", doc.OriginalName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentPostgres_List_WithoutStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	userID := 1
	now := time.Now().Round(time.Second)

	rows := sqlmock.NewRows([]string{"id", "user_id", "original_name", "s3_key", "mime_type", "file_size", "status", "error_msg", "created_at", "updated_at"}).
		AddRow(1, 1, "doc1.pdf", "key1", "application/pdf", 100, "completed", "", now, now).
		AddRow(2, 1, "doc2.pdf", "key2", "application/pdf", 200, "processing", "", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM documents WHERE user_id = $1 ORDER BY created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(rows)

	docs, err := repo.List(userID, "")
	assert.NoError(t, err)
	assert.Len(t, docs, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentPostgres_List_WithStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	userID := 1
	status := "completed"
	now := time.Now().Round(time.Second)

	rows := sqlmock.NewRows([]string{"id", "user_id", "original_name", "s3_key", "mime_type", "file_size", "status", "error_msg", "created_at", "updated_at"}).
		AddRow(1, 1, "doc1.pdf", "key1", "application/pdf", 100, "completed", "", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM documents WHERE user_id = $1 AND status = $2 ORDER BY created_at DESC`)).
		WithArgs(userID, status).
		WillReturnRows(rows)

	docs, err := repo.List(userID, status)
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "completed", docs[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentPostgres_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	docID := 10
	status := "completed"
	errorMsg := ""

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE documents SET status = $1, error_msg = $2, updated_at = NOW() WHERE id = $3`)).
		WithArgs(status, errorMsg, docID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateStatus(docID, status, errorMsg)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentPostgres_CheckBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	userID := 1
	expectedBalance := 15.75

	rows := sqlmock.NewRows([]string{"balance"}).AddRow(expectedBalance)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT balance FROM user_balances WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	balance, err := repo.CheckBalance(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentPostgres_DecrementBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDocumentPostgres(sqlxDB)

	userID := 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE user_balances SET balance = balance - 1, updated_at = NOW() WHERE user_id = $1 AND balance >= 1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DecrementBalance(userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
