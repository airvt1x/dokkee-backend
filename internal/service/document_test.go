package service

import (
	"errors"
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Моки для репозиториев

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Create(doc dokkee.Document) (int, error) {
	args := m.Called(doc)
	return args.Int(0), args.Error(1)
}

func (m *MockDocumentRepository) GetByID(docID, userID int) (dokkee.Document, error) {
	args := m.Called(docID, userID)
	return args.Get(0).(dokkee.Document), args.Error(1)
}

func (m *MockDocumentRepository) List(userID int, status string) ([]dokkee.Document, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]dokkee.Document), args.Error(1)
}

func (m *MockDocumentRepository) UpdateStatus(docID int, status, errorMsg string) error {
	args := m.Called(docID, status, errorMsg)
	return args.Error(0)
}

func (m *MockDocumentRepository) CheckBalance(userID int) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDocumentRepository) DecrementBalance(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

type MockResultRepository struct {
	mock.Mock
}

func (m *MockResultRepository) Save(result dokkee.AnalysisResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockResultRepository) GetByDocumentID(docID int) (dokkee.AnalysisResult, error) {
	args := m.Called(docID)
	return args.Get(0).(dokkee.AnalysisResult), args.Error(1)
}

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Upload(key string, file interface{}, contentType string) error {
	args := m.Called(key, file, contentType)
	return args.Error(0)
}

func (m *MockFileStorage) Download(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

// Тесты для DocumentService

func TestDocumentService_GetByID_Success(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{
		repo: mockDocRepo,
	}

	expected := dokkee.Document{Id: 123, UserID: 1, OriginalName: "test.pdf"}
	mockDocRepo.On("GetByID", 123, 1).Return(expected, nil)

	doc, err := svc.GetByID(123, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, doc)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_GetByID_NotFound(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{
		repo: mockDocRepo,
	}

	mockDocRepo.On("GetByID", 999, 1).Return(dokkee.Document{}, errors.New("not found"))

	_, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Equal(t, "document not found", err.Error())
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_List_Empty(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{
		repo: mockDocRepo,
	}

	mockDocRepo.On("List", 1, "").Return([]dokkee.Document{}, nil)

	docs, err := svc.List(1, "")
	assert.NoError(t, err)
	assert.Len(t, docs, 0)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_List_WithResults(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{
		repo: mockDocRepo,
	}

	expected := []dokkee.Document{
		{Id: 1, OriginalName: "doc1.pdf"},
		{Id: 2, OriginalName: "doc2.pdf"},
	}
	mockDocRepo.On("List", 1, "").Return(expected, nil)

	docs, err := svc.List(1, "")
	assert.NoError(t, err)
	assert.Len(t, docs, 2)
	assert.Equal(t, expected, docs)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_List_WithStatus(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{
		repo: mockDocRepo,
	}

	expected := []dokkee.Document{{Id: 1, Status: "completed"}}
	mockDocRepo.On("List", 1, "completed").Return(expected, nil)

	docs, err := svc.List(1, "completed")
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "completed", docs[0].Status)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_List_Error(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{
		repo: mockDocRepo,
	}

	mockDocRepo.On("List", 1, "").Return([]dokkee.Document{}, errors.New("db error"))

	_, err := svc.List(1, "")
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_GetByID_EmptyResult(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{repo: mockDocRepo}

	mockDocRepo.On("GetByID", 999, 1).Return(dokkee.Document{}, nil)

	doc, err := svc.GetByID(999, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, doc.Id)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_List_EmptyUserID(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{repo: mockDocRepo}

	mockDocRepo.On("List", 0, "").Return([]dokkee.Document{}, nil)

	docs, err := svc.List(0, "")
	assert.NoError(t, err)
	assert.Len(t, docs, 0)
	mockDocRepo.AssertExpectations(t)
}


func TestDocumentService_CheckBalance_Error(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{repo: mockDocRepo}

	mockDocRepo.On("CheckBalance", 1).Return(0.0, errors.New("db connection error"))

	err := svc.checkBalance(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check balance")
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_CheckBalance_Success(t *testing.T) {
	mockDocRepo := new(MockDocumentRepository)
	svc := &DocumentService{repo: mockDocRepo}

	mockDocRepo.On("CheckBalance", 1).Return(5.0, nil)

	err := svc.checkBalance(1)
	assert.NoError(t, err)
	mockDocRepo.AssertExpectations(t)
}