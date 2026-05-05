package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) uploadDocument(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user not found in context")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	docID, err := h.services.Document.Upload(userID, file, header)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"id": docID})
}

func (h *Handler) listDocuments(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user not found in context")
		return
	}

	status := c.Query("status")
	docs, err := h.services.Document.List(userID, status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

func (h *Handler) getDocument(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user not found in context")
		return
	}

	docID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid document id")
		return
	}

	doc, err := h.services.Document.GetByID(docID, userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, doc)
}
