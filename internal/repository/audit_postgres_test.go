package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuditPostgres_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuditPostgres(sqlxDB)

	record := dokkee.AuditRecord{
		EventType:  "DOCUMENT_UPLOADED",
		UserIDHash: "abc123",
		DocIDHash:  nil,
		IPHash:     "iphash",
		Success:    true,
		ErrorCode:  "",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO audit_log (event_type, user_id_hash, doc_id_hash, ip_hash, success, error_code) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(record.EventType, record.UserIDHash, record.DocIDHash, record.IPHash, record.Success, record.ErrorCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(record)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditPostgres_Insert_WithDocIDHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuditPostgres(sqlxDB)

	docHash := "doc123hash"
	record := dokkee.AuditRecord{
		EventType:  "RESULT_ACCESSED",
		UserIDHash: "user456",
		DocIDHash:  &docHash,
		IPHash:     "ip789",
		Success:    true,
		ErrorCode:  "",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO audit_log (event_type, user_id_hash, doc_id_hash, ip_hash, success, error_code) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(record.EventType, record.UserIDHash, record.DocIDHash, record.IPHash, record.Success, record.ErrorCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(record)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditPostgres_Insert_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAuditPostgres(sqlxDB)

	record := dokkee.AuditRecord{
		EventType:  "API_REQUEST",
		UserIDHash: "hash123",
		IPHash:     "iphash",
		Success:    true,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO audit_log (event_type, user_id_hash, doc_id_hash, ip_hash, success, error_code) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(record.EventType, record.UserIDHash, record.DocIDHash, record.IPHash, record.Success, record.ErrorCode).
		WillReturnError(errors.New("db error"))

	err = repo.Insert(record)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
