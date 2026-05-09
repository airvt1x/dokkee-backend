package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResultService struct {
	mock.Mock
}

func (m *MockResultService) GetByDocumentID(docID int) (dokkee.AnalysisResult, error) {
	args := m.Called(docID)
	return args.Get(0).(dokkee.AnalysisResult), args.Error(1)
}

func TestHandler_getResult_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	mockResult := new(MockResultService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
			Result:   mockResult,
		},
	}

	doc := dokkee.Document{Id: 10, Status: "completed"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)

	result := dokkee.AnalysisResult{DocumentID: 10, ResultJSON: []byte(`{"risk":"high"}`)}
	mockResult.On("GetByDocumentID", 10).Return(result, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDoc.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestHandler_getResult_Processing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	doc := dokkee.Document{Id: 10, Status: "processing"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "processing", resp["status"])
	mockDoc.AssertExpectations(t)
}

func TestHandler_getResult_DocumentNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	mockDoc.On("GetByID", 99, 1).Return(dokkee.Document{}, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/99/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDoc.AssertExpectations(t)
}

func TestHandler_getResult_ResultNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	mockResult := new(MockResultService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
			Result:   mockResult,
		},
	}

	doc := dokkee.Document{Id: 10, Status: "completed"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)
	mockResult.On("GetByDocumentID", 10).Return(dokkee.AnalysisResult{}, errors.New("result not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDoc.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestHandler_getResult_WrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDoc := new(MockDocumentService)
	handler := &Handler{
		services: &service.Service{
			Document: mockDoc,
		},
	}

	mockDoc.On("GetByID", 10, 1).Return(dokkee.Document{}, errors.New("document not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockDoc.AssertExpectations(t)
}

func TestHandler_getResult_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/abc/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_getResult_EmptyUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	// Не устанавливаем user_id
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
}

func TestHandler_getResult_EmptyDocID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodGet, "/api/documents//result", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/documents/:id/result", handler.getResult)
	router.ServeHTTP(w, req)

	// Пустой ID приводит к 400 Bad Request (invalid document id)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}