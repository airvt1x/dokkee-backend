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
)

func TestHandler_getProfile_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	mockAuth.On("GetProfile", 1).Return(dokkee.User{}, errors.New("user not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.GET("/api/profile", handler.getProfile)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockAuth.AssertExpectations(t)
}

func TestHandler_updateProfile_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{services: &service.Service{}}

	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.PATCH("/api/profile", handler.updateProfile)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_updateProfile_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := new(MockAuthorizationService)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	input := dokkee.UpdateProfileInput{}
	mockAuth.On("UpdateProfile", 1, input).Return(errors.New("update failed"))

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPatch, "/api/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	router.PATCH("/api/profile", handler.updateProfile)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockAuth.AssertExpectations(t)
}