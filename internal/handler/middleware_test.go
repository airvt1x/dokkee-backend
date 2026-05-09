package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/airvt1x/dokkee-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) Log(event service.AuditEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestJwtMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	token := "valid.token"
	mockAuth.On("ParseToken", token).Return(1, nil)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.jwtMiddleware()(c)

	assert.Equal(t, http.StatusOK, w.Code)
	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, 1, userID)
	mockAuth.AssertExpectations(t)
}

func TestJwtMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.jwtMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authorization header is empty")
}

func TestJwtMiddleware_WrongFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidScheme token")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.jwtMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
}

func TestJwtMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	token := "invalid.token"
	mockAuth.On("ParseToken", token).Return(0, assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.jwtMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuth.AssertExpectations(t)
}

func TestGetUserID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 42)

	userID, ok := getUserID(c)
	assert.True(t, ok)
	assert.Equal(t, 42, userID)
}

func TestGetUserID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID, ok := getUserID(c)
	assert.False(t, ok)
	assert.Equal(t, 0, userID)
}

func TestAuditMiddleware_CallsLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Audit: mockAudit,
		},
	}

	mockAudit.On("Log", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/documents", nil)
	c.Set("user_id", 1)

	handler.auditMiddleware()(c)

	mockAudit.AssertExpectations(t)
}

func TestAuditEventByPath(t *testing.T) {
	tests := []struct {
		path   string
		method string
		want   string
	}{
		{"/api/documents", "POST", "DOCUMENT_UPLOADED"},
		{"/api/documents/1/result", "GET", "RESULT_ACCESSED"},
		{"/api/profile", "GET", "PROFILE_ACCESSED"},
		{"/api/profile", "PATCH", "PROFILE_UPDATED"},
		{"/api/unknown", "GET", "API_REQUEST"},
		{"/", "GET", "API_REQUEST"},
	}

	for _, tt := range tests {
		t.Run(tt.path+"_"+tt.method, func(t *testing.T) {
			got := auditEventByPath(tt.path, tt.method)
			assert.Equal(t, tt.want, got)
		})
	}
}