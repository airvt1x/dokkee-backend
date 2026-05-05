package handler

import (
	"github.com/airvt1x/dokkee-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.jwtMiddleware(), h.auditMiddleware())
	{
		docs := api.Group("/documents")
		{
			docs.POST("", h.uploadDocument)
			docs.GET("", h.listDocuments)
			docs.GET("/:id", h.getDocument)
			docs.GET("/:id/result", h.getResult)
		}

		api.GET("/profile", h.getProfile)
		api.PATCH("/profile", h.updateProfile)
	}

	return router
}
