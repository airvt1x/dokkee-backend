package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) Upload(userID int, file multipart.File, header *multipart.FileHeader) (int, error) {
	args := m.Called(userID, file, header)
	return args.Int(0), args.Error(1)
}

func (m *MockDocumentService) GetByID(docID, userID int) (dokkee.Document, error) {
	args := m.Called(docID, userID)
	return args.Get(0).(dokkee.Document), args.Error(1)
}

func (m *MockDocumentService) List(userID int, status string) ([]dokkee.Document, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]dokkee.Document), args.Error(1)
}

func createMultipartFile(t *testing.T, content, filename string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)
	_, err = part.Write([]byte(content))
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)
	return body, writer.FormDataContentType()
}

// ========== СУЩЕСТВУЮЩИЕ ТЕСТЫ ==========

func TestHandler_uploadDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	body, contentType := createMultipartFile(t, "test content", "test.pdf")
	mockDoc.On("Upload", 1, mock.Anything, mock.Anything).Return(123, nil)

	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.POST("/api/documents", handler.uploadDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), resp["id"])
	mockDoc.AssertExpectations(t)
}

func TestHandler_uploadDocument_NoFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodPost, "/api/documents", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.POST("/api/documents", handler.uploadDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_listDocuments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	docs := []dokkee.Document{{Id: 1, OriginalName: "doc1.pdf"}, {Id: 2, OriginalName: "doc2.pdf"}}
	mockDoc.On("List", 1, "").Return(docs, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents", handler.listDocuments)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp["documents"], 2)
	mockDoc.AssertExpectations(t)
}

func TestHandler_listDocuments_WithStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	status := "completed"
	docs := []dokkee.Document{{Id: 1, Status: "completed"}}
	mockDoc.On("List", 1, status).Return(docs, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents?status=completed", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents", handler.listDocuments)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDoc.AssertExpectations(t)
}

func TestHandler_getDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	doc := dokkee.Document{Id: 5, OriginalName: "test.pdf"}
	mockDoc.On("GetByID", 5, 1).Return(doc, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/5", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id", handler.getDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dokkee.Document
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, doc, resp)
	mockDoc.AssertExpectations(t)
}

func TestHandler_getDocument_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	mockDoc.On("GetByID", 99, 1).Return(dokkee.Document{}, errors.New("document not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/99", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id", handler.getDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDoc.AssertExpectations(t)
}

// ========== ДОПОЛНИТЕЛЬНЫЕ ТЕСТЫ ==========

func TestHandler_uploadDocument_InsufficientBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	body, contentType := createMultipartFile(t, "test content", "test.pdf")
	mockDoc.On("Upload", 1, mock.Anything, mock.Anything).Return(0, errors.New("insufficient balance"))

	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.POST("/api/documents", handler.uploadDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockDoc.AssertExpectations(t)
}

func TestHandler_listDocuments_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	mockDoc.On("List", 1, "").Return([]dokkee.Document{}, errors.New("db error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents", handler.listDocuments)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDoc.AssertExpectations(t)
}

func TestHandler_getDocument_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/abc", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id", handler.getDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_getDocument_WrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	mockDoc.On("GetByID", 5, 1).Return(dokkee.Document{}, errors.New("document not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/5", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id", handler.getDocument)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDoc.AssertExpectations(t)
}
