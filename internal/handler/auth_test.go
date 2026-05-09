package handler

import (
	"bytes"
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

type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) CreateUser(user dokkee.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthorizationService) GenerateToken(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorizationService) ParseToken(token string) (int, error) {
	args := m.Called(token)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthorizationService) GetProfile(userID int) (dokkee.User, error) {
	args := m.Called(userID)
	return args.Get(0).(dokkee.User), args.Error(1)
}

func (m *MockAuthorizationService) UpdateProfile(userID int, input dokkee.UpdateProfileInput) error {
	args := m.Called(userID, input)
	return args.Error(0)
}

// ========== СУЩЕСТВУЮЩИЕ ТЕСТЫ ==========

func TestHandler_signUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockService,
		},
	}

	user := dokkee.User{
		Username:  "testuser",
		Password:  "password",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Phone:     "+1234567890",
	}

	mockService.On("CreateUser", mock.AnythingOfType("dokkee.User")).Return(1, nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/auth/sign-up", handler.signUp)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["id"])
	mockService.AssertExpectations(t)
}

func TestHandler_signIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockService,
		},
	}

	input := signInInput{
		Username: "testuser",
		Password: "password",
	}

	mockService.On("GenerateToken", input.Username, input.Password).Return("token123", nil)

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/auth/sign-in", handler.signIn)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "token123", response["token"])
	mockService.AssertExpectations(t)
}

// ========== НОВЫЕ ТЕСТЫ ==========

func TestHandler_signUp_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/auth/sign-up", handler.signUp)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_signUp_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockService,
		},
	}

	user := dokkee.User{
		Username:  "testuser",
		Password:  "password",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Phone:     "+1234567890",
	}
	mockService.On("CreateUser", mock.Anything).Return(0, errors.New("duplicate username"))

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/auth/sign-up", handler.signUp)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_signIn_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/auth/sign-in", handler.signIn)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
