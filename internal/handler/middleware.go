package handler

import (
	"net/http"
	"strings"

	"github.com/airvt1x/dokkee-backend/internal/service"
	"github.com/gin-gonic/gin"
)

const userCtx = "user_id"

func (h *Handler) jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is empty"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		userID, err := h.services.Authorization.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(userCtx, userID)
		c.Next()
	}
}

func (h *Handler) auditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		userID, exists := c.Get(userCtx)
		if !exists {
			return
		}

		status := c.Writer.Status()
		success := status >= 200 && status < 400

		h.services.Audit.Log(service.AuditEvent{
			Type:    auditEventByPath(c.FullPath(), c.Request.Method),
			UserID:  userID.(int),
			IP:      c.ClientIP(),
			Success: success,
		})
	}
}

func auditEventByPath(path, method string) string {
	switch {
	case strings.Contains(path, "/documents") && method == "POST":
		return "DOCUMENT_UPLOADED"
	case strings.Contains(path, "/result"):
		return "RESULT_ACCESSED"
	case strings.Contains(path, "/profile") && method == "GET":
		return "PROFILE_ACCESSED"
	case strings.Contains(path, "/profile") && method == "PATCH":
		return "PROFILE_UPDATED"
	default:
		return "API_REQUEST"
	}
}

func getUserID(c *gin.Context) (int, bool) {
	id, exists := c.Get(userCtx)
	if !exists {
		return 0, false
	}
	userID, ok := id.(int)
	return userID, ok
}
