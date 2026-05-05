package repository

import (
	"mime/multipart"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user dokkee.User) (int, error)
	GetUser(username string) (dokkee.User, error)
	GetProfile(userID int) (dokkee.User, error)
	UpdateProfile(userID int, input dokkee.UpdateProfileInput) error
}

type Document interface {
	Create(doc dokkee.Document) (int, error)
	GetByID(docID, userID int) (dokkee.Document, error)
	List(userID int, status string) ([]dokkee.Document, error)
	UpdateStatus(docID int, status, errorMsg string) error
	CheckBalance(userID int) (float64, error)
	DecrementBalance(userID int) error
}

type Result interface {
	Save(result dokkee.AnalysisResult) error
	GetByDocumentID(docID int) (dokkee.AnalysisResult, error)
}

type Audit interface {
	Insert(record dokkee.AuditRecord) error
}

type FileStorage interface {
	Upload(key string, file multipart.File, contentType string) error
	Download(key string) ([]byte, error)
}

type Repository struct {
	Authorization
	Document
	Result
	Audit
	FileStorage
}

func NewRepository(db *sqlx.DB, s3 *S3Repository) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Document:      NewDocumentPostgres(db),
		Result:        NewResultPostgres(db),
		Audit:         NewAuditPostgres(db),
		FileStorage:   s3,
	}
}
