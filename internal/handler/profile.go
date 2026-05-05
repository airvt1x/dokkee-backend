package handler

import (
	"net/http"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user not found in context")
		return
	}

	user, err := h.services.Authorization.GetProfile(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user not found in context")
		return
	}

	var input dokkee.UpdateProfileInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Authorization.UpdateProfile(userID, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
