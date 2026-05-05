package dokkee

import "time"

type Document struct {
	Id           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	OriginalName string    `json:"original_name" db:"original_name"`
	S3Key        string    `json:"s3_key" db:"s3_key"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	Status       string    `json:"status" db:"status"`
	ErrorMsg     string    `json:"error_msg,omitempty" db:"error_msg"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type AnalysisResult struct {
	Id         int       `json:"id" db:"id"`
	DocumentID int       `json:"document_id" db:"document_id"`
	ResultJSON []byte    `json:"result" db:"result_json"`
	ModelUsed  string    `json:"model_used" db:"model_used"`
	TokensUsed int       `json:"tokens_used" db:"tokens_used"`
	AnalyzedAt time.Time `json:"analyzed_at" db:"analyzed_at"`
}

type AuditRecord struct {
	EventType  string  `db:"event_type"`
	UserIDHash string  `db:"user_id_hash"`
	DocIDHash  *string `db:"doc_id_hash"`
	IPHash     string  `db:"ip_hash"`
	Success    bool    `db:"success"`
	ErrorCode  string  `db:"error_code"`
}
