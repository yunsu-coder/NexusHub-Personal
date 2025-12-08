package handler

import (
	"net/http"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CodeHandler struct {
	service *service.CodeService
}

func NewCodeHandler() *CodeHandler {
	return &CodeHandler{
		service: service.NewCodeService(),
	}
}

func (h *CodeHandler) GetAll(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	snippets, err := h.service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, snippets)
}

func (h *CodeHandler) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	snippet, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code snippet not found"})
		return
	}
	c.JSON(http.StatusOK, snippet)
}

func (h *CodeHandler) GetByLanguage(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	language := c.Param("language")

	snippets, err := h.service.GetByLanguage(language, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, snippets)
}

func (h *CodeHandler) Create(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var snippet model.CodeSnippet
	if err := c.ShouldBindJSON(&snippet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippet.UserID = userID
	if err := h.service.Create(&snippet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, snippet)
}

func (h *CodeHandler) Update(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var snippet model.CodeSnippet
	if err := c.ShouldBindJSON(&snippet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippet.ID = uint(id)
	snippet.UserID = userID
	if err := h.service.Update(&snippet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, snippet)
}

func (h *CodeHandler) Delete(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Code snippet deleted successfully"})
}
