package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/airvt1x/dokkee-backend"
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

func (m *MockAuthorizationService) GenerateToken(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorizationService) ParseToken(token string) (int, error) {
	args := m.Called(token)
	return args.Int(0), args.Error(1)
}

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
		Email:    "test@example.com",
		Password: "password",
	}

	mockService.On("GenerateToken", input.Email, input.Password).Return("token123", nil)

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
