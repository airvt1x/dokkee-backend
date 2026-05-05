package service

import (
	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/repository"
)

type ResultService struct {
	repo repository.Result
}

func NewResultService(repo repository.Result) *ResultService {
	return &ResultService{repo: repo}
}

func (s *ResultService) GetByDocumentID(docID int) (dokkee.AnalysisResult, error) {
	return s.repo.GetByDocumentID(docID)
}
