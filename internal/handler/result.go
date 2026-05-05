package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getResult(c *gin.Context) {
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

	// Проверка владельца документа
	doc, err := h.services.Document.GetByID(docID, userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "document not found")
		return
	}

	if doc.Status != "completed" {
		c.JSON(http.StatusAccepted, gin.H{
			"status":  doc.Status,
			"message": "analysis is not yet complete",
		})
		return
	}

	result, err := h.services.Result.GetByDocumentID(docID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "result not found")
		return
	}

	c.JSON(http.StatusOK, result)
}
