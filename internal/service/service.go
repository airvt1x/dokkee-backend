package service

import (
	"mime/multipart"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/repository"
)

type Authorization interface {
	CreateUser(user dokkee.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GetProfile(userID int) (dokkee.User, error)
	UpdateProfile(userID int, input dokkee.UpdateProfileInput) error
}

type Documents interface {
	Upload(userID int, file multipart.File, header *multipart.FileHeader) (int, error)
	GetByID(docID, userID int) (dokkee.Document, error)
	List(userID int, status string) ([]dokkee.Document, error)
}

type Results interface {
	GetByDocumentID(docID int) (dokkee.AnalysisResult, error)
}

type Audit interface {
	Log(event AuditEvent) error
}

type Service struct {
	Authorization
	Document Documents
	Result   Results
	Audit    Audit
}

func NewService(repos *repository.Repository) *Service {
	ner := NewNERService()
	aiConn := NewAIConnector()

	docService := NewDocumentService(repos.Document, repos.Result, repos.FileStorage, ner, aiConn)

	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Document:      docService,
		Result:        NewResultService(repos.Result),
		Audit:         NewAuditService(repos.Audit),
	}
}
