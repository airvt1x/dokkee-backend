package repository

import (
	"fmt"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
)

type DocumentPostgres struct {
	db *sqlx.DB
}

func NewDocumentPostgres(db *sqlx.DB) *DocumentPostgres {
	return &DocumentPostgres{db: db}
}

func (r *DocumentPostgres) Create(doc dokkee.Document) (int, error) {
	var id int
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, original_name, s3_key, mime_type, file_size, status)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, documentsTable)
	err := r.db.QueryRow(query,
		doc.UserID, doc.OriginalName, doc.S3Key, doc.MimeType, doc.FileSize, doc.Status,
	).Scan(&id)
	return id, err
}

func (r *DocumentPostgres) GetByID(docID, userID int) (dokkee.Document, error) {
	var doc dokkee.Document
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1 AND user_id = $2`, documentsTable)
	err := r.db.Get(&doc, query, docID, userID)
	return doc, err
}

func (r *DocumentPostgres) List(userID int, status string) ([]dokkee.Document, error) {
	var docs []dokkee.Document
	var query string
	var args []interface{}

	if status != "" {
		query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 AND status = $2 ORDER BY created_at DESC`, documentsTable)
		args = []interface{}{userID, status}
	} else {
		query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 ORDER BY created_at DESC`, documentsTable)
		args = []interface{}{userID}
	}

	err := r.db.Select(&docs, query, args...)
	return docs, err
}

func (r *DocumentPostgres) UpdateStatus(docID int, status, errorMsg string) error {
	query := fmt.Sprintf(`
		UPDATE %s SET status = $1, error_msg = $2, updated_at = NOW() WHERE id = $3`, documentsTable)
	_, err := r.db.Exec(query, status, errorMsg, docID)
	return err
}

func (r *DocumentPostgres) CheckBalance(userID int) (float64, error) {
	var balance float64
	query := fmt.Sprintf(`SELECT balance FROM %s WHERE user_id = $1`, userBalancesTable)
	err := r.db.QueryRow(query, userID).Scan(&balance)
	return balance, err
}

func (r *DocumentPostgres) DecrementBalance(userID int) error {
	query := fmt.Sprintf(`UPDATE %s SET balance = balance - 1, updated_at = NOW() WHERE user_id = $1 AND balance >= 1`, userBalancesTable)
	_, err := r.db.Exec(query, userID)
	return err
}
