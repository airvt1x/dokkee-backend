package repository

import (
	"fmt"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
)

type AuditPostgres struct {
	db *sqlx.DB
}

func NewAuditPostgres(db *sqlx.DB) *AuditPostgres {
	return &AuditPostgres{db: db}
}

func (r *AuditPostgres) Insert(record dokkee.AuditRecord) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (event_type, user_id_hash, doc_id_hash, ip_hash, success, error_code)
		VALUES ($1, $2, $3, $4, $5, $6)`, auditLogTable)
	_, err := r.db.Exec(query,
		record.EventType, record.UserIDHash, record.DocIDHash,
		record.IPHash, record.Success, record.ErrorCode,
	)
	return err
}
