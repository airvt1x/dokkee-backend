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

// Helper functions

func setupRouterWithMocks(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return h.InitRoutes()
}

func createMultipartFileBuffer(t *testing.T, content, filename string) (*bytes.Buffer, string) {
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

// ========== AUTH ENDPOINTS ==========

func TestAPI_SignUp_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	user := dokkee.User{
		Username:  "apiuser",
		Password:  "pass123",
		FirstName: "API",
		LastName:  "User",
		Email:     "api@test.com",
		Phone:     "+79991112299",
	}
	mockAuth.On("CreateUser", mock.AnythingOfType("dokkee.User")).Return(42, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(42), resp["id"])
	mockAuth.AssertExpectations(t)
}

func TestAPI_SignUp_InvalidJSON(t *testing.T) {
	handler := &Handler{services: &service.Service{}}
	router := setupRouterWithMocks(handler)

	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_SignUp_MissingFields(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	handler := &Handler{services: &service.Service{Authorization: mockAuth}}
	router := setupRouterWithMocks(handler)

	user := map[string]string{"username": "only_username"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_SignIn_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	input := signInInput{Username: "apiuser", Password: "pass123"}
	mockAuth.On("GenerateToken", input.Username, input.Password).Return("valid-token", nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "valid-token", resp["token"])
	mockAuth.AssertExpectations(t)
}

func TestAPI_SignIn_InvalidCredentials(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	input := signInInput{Username: "wrong", Password: "wrong"}
	mockAuth.On("GenerateToken", input.Username, input.Password).Return("", errors.New("invalid credentials"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuth.AssertExpectations(t)
}

func TestAPI_SignIn_InvalidJSON(t *testing.T) {
	handler := &Handler{services: &service.Service{}}
	router := setupRouterWithMocks(handler)

	req, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ========== DOCUMENTS ENDPOINTS ==========

func TestAPI_UploadDocument_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockDoc.On("Upload", 1, mock.Anything, mock.Anything).Return(100, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, contentType := createMultipartFileBuffer(t, "test content", "test.pdf")
	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(100), resp["id"])
}

func TestAPI_UploadDocument_InvalidFileType(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockDoc.On("Upload", 1, mock.Anything, mock.Anything).Return(0, errors.New("unsupported file type"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, contentType := createMultipartFileBuffer(t, "test content", "test.exe")
	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_UploadDocument_Unauthorized(t *testing.T) {
	handler := &Handler{services: &service.Service{}}
	router := setupRouterWithMocks(handler)

	body, contentType := createMultipartFileBuffer(t, "test", "test.pdf")
	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPI_UploadDocument_NoFile(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodPost, "/api/documents", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_ListDocuments_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	docs := []dokkee.Document{{Id: 1, OriginalName: "doc1.pdf"}, {Id: 2, OriginalName: "doc2.pdf"}}
	mockDoc.On("List", 1, "").Return(docs, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp["documents"], 2)
}

func TestAPI_ListDocuments_WithStatus(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	docs := []dokkee.Document{{Id: 1, Status: "completed"}}
	mockDoc.On("List", 1, "completed").Return(docs, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents?status=completed", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_ListDocuments_ServiceError(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockDoc.On("List", 1, "").Return([]dokkee.Document{}, errors.New("database error"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAPI_GetDocument_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	doc := dokkee.Document{Id: 5, OriginalName: "test.pdf"}
	mockDoc.On("GetByID", 5, 1).Return(doc, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/5", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetDocument_NotFound(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockDoc.On("GetByID", 99, 1).Return(dokkee.Document{}, errors.New("not found"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/99", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_GetDocument_InvalidID(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/abc", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ========== RESULT ENDPOINT ==========

func TestAPI_GetResult_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockResult := new(MockResultService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Result:        mockResult,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	doc := dokkee.Document{Id: 10, Status: "completed"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)
	result := dokkee.AnalysisResult{DocumentID: 10, ResultJSON: []byte(`{"risk":"low"}`)}
	mockResult.On("GetByDocumentID", 10).Return(result, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetResult_Processing(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	doc := dokkee.Document{Id: 10, Status: "processing"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "processing", resp["status"])
}

func TestAPI_GetResult_Failed(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	doc := dokkee.Document{Id: 10, Status: "failed", ErrorMsg: "analysis error"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "failed", resp["status"])
	assert.Equal(t, "analysis is not yet complete", resp["message"])
}

func TestAPI_GetResult_InvalidID(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/abc/result", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ========== PROFILE ENDPOINTS ==========

func TestAPI_GetProfile_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	user := dokkee.User{Id: 1, Username: "testuser", FirstName: "Test"}
	mockAuth.On("GetProfile", 1).Return(user, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetProfile_Unauthorized(t *testing.T) {
	handler := &Handler{services: &service.Service{}}
	router := setupRouterWithMocks(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPI_UpdateProfile_Success(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAuth.On("UpdateProfile", 1, mock.Anything).Return(nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	input := map[string]string{"first_name": "UpdatedName"}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_UpdateProfile_InvalidJSON(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBufferString("{invalid}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_UpdateProfile_EmptyBody(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAuth.On("UpdateProfile", 1, mock.Anything).Return(nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ========== MIDDLEWARE TESTS ==========

func TestAPI_JwtMiddleware_NoToken(t *testing.T) {
	handler := &Handler{services: &service.Service{}}
	router := setupRouterWithMocks(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPI_JwtMiddleware_InvalidToken(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	handler := &Handler{services: &service.Service{Authorization: mockAuth}}

	mockAuth.On("ParseToken", "invalid").Return(0, errors.New("token error"))

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPI_JwtMiddleware_WrongFormat(t *testing.T) {
	handler := &Handler{services: &service.Service{}}
	router := setupRouterWithMocks(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/documents", nil)
	req.Header.Set("Authorization", "WrongFormat token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ========== AUDIT MIDDLEWARE ==========

type MockAuditServiceForAPI struct {
	mock.Mock
}

func (m *MockAuditServiceForAPI) Log(event service.AuditEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestAPI_UploadDocument_TooLarge(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockDoc.On("Upload", 1, mock.Anything, mock.Anything).Return(0, errors.New("file too large: max 10 MB"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, contentType := createMultipartFileBuffer(t, string(make([]byte, 11*1024*1024)), "large.pdf")
	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_GetDocument_WrongUser(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "token1").Return(1, nil)
	mockDoc.On("GetByID", 5, 1).Return(dokkee.Document{}, errors.New("document not found"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/5", nil)
	req.Header.Set("Authorization", "Bearer token1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_UpdateProfile_PartialUpdate(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAuth.On("UpdateProfile", 1, mock.MatchedBy(func(input dokkee.UpdateProfileInput) bool {
		return input.FirstName != nil && *input.FirstName == "JustName" && input.LastName == nil
	})).Return(nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	input := map[string]string{"first_name": "JustName"}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetResult_Queued(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	doc := dokkee.Document{Id: 10, Status: "queued"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "queued", resp["status"])
}

func TestAPI_UploadDocument_UnsupportedMimeType(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockDoc.On("Upload", 1, mock.Anything, mock.Anything).Return(0, errors.New("unsupported file type"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, _ := createMultipartFileBuffer(t, "test", "test.exe")
	req, _ := http.NewRequest(http.MethodPost, "/api/documents", body)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_GetResult_NotFound(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockDoc := new(MockDocumentService)
	mockResult := new(MockResultService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Document:      mockDoc,
			Result:        mockResult,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	doc := dokkee.Document{Id: 10, Status: "completed"}
	mockDoc.On("GetByID", 10, 1).Return(doc, nil)
	mockResult.On("GetByDocumentID", 10).Return(dokkee.AnalysisResult{}, errors.New("result not found"))
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	req, _ := http.NewRequest(http.MethodGet, "/api/documents/10/result", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_UpdateProfile_NothingToUpdate(t *testing.T) {
	mockAuth := new(MockAuthorizationService)
	mockAudit := new(MockAuditService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Audit:         mockAudit,
		},
	}

	mockAuth.On("ParseToken", "valid-token").Return(1, nil)
	mockAuth.On("UpdateProfile", 1, mock.Anything).Return(nil)
	mockAudit.On("Log", mock.Anything).Return(nil)

	router := setupRouterWithMocks(handler)
	body, _ := json.Marshal(map[string]interface{}{})
	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}