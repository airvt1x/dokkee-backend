package service

import (
	"errors"
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResultRepositoryForService struct {
	mock.Mock
}

func (m *MockResultRepositoryForService) Save(result dokkee.AnalysisResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockResultRepositoryForService) GetByDocumentID(docID int) (dokkee.AnalysisResult, error) {
	args := m.Called(docID)
	return args.Get(0).(dokkee.AnalysisResult), args.Error(1)
}

func TestResultService_GetByDocumentID_Success(t *testing.T) {
	mockRepo := new(MockResultRepositoryForService)
	svc := NewResultService(mockRepo)

	expected := dokkee.AnalysisResult{DocumentID: 10, ResultJSON: []byte(`{"risk":"high"}`)}
	mockRepo.On("GetByDocumentID", 10).Return(expected, nil)

	result, err := svc.GetByDocumentID(10)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestResultService_GetByDocumentID_NotFound(t *testing.T) {
	mockRepo := new(MockResultRepositoryForService)
	svc := NewResultService(mockRepo)

	mockRepo.On("GetByDocumentID", 999).Return(dokkee.AnalysisResult{}, errors.New("result not found"))

	_, err := svc.GetByDocumentID(999)
	assert.Error(t, err)
	assert.Equal(t, "result not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestResultService_GetByDocumentID_EmptyResult(t *testing.T) {
	mockRepo := new(MockResultRepositoryForService)
	svc := NewResultService(mockRepo)

	mockRepo.On("GetByDocumentID", 1).Return(dokkee.AnalysisResult{}, nil)

	result, err := svc.GetByDocumentID(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.DocumentID)
	assert.Nil(t, result.ResultJSON)
	mockRepo.AssertExpectations(t)
}
