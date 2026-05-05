package repository

import (
	"fmt"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
)

type ResultPostgres struct {
	db *sqlx.DB
}

func NewResultPostgres(db *sqlx.DB) *ResultPostgres {
	return &ResultPostgres{db: db}
}

func (r *ResultPostgres) Save(result dokkee.AnalysisResult) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (document_id, result_json, model_used, tokens_used)
		VALUES ($1, $2, $3, $4)`, analysisResultsTable)
	_, err := r.db.Exec(query, result.DocumentID, result.ResultJSON, result.ModelUsed, result.TokensUsed)
	return err
}

func (r *ResultPostgres) GetByDocumentID(docID int) (dokkee.AnalysisResult, error) {
	var result dokkee.AnalysisResult
	query := fmt.Sprintf(`SELECT * FROM %s WHERE document_id = $1`, analysisResultsTable)
	err := r.db.Get(&result, query, docID)
	return result, err
}
