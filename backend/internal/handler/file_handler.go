package handler

import (
	"net/http"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *service.FileService
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		service: service.NewFileService(),
	}
}

func (h *FileHandler) GetAll(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	files, err := h.service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	file, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}

func (h *FileHandler) GetByCategory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	category := c.Param("category")

	files, err := h.service.GetByCategory(category, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) Upload(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	uploadedFile, err := h.service.Upload(file, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate accessible URL path for the uploaded file
	// Extract the relative path from uploads/ onwards
	relativePath := ""
	if idx := strings.Index(uploadedFile.FilePath, "uploads"); idx != -1 {
		relativePath = uploadedFile.FilePath[idx:]
		// Convert backslashes to forward slashes for URL
		relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         uploadedFile.ID,
		"file_name":  uploadedFile.FileName,
		"file_size":  uploadedFile.FileSize,
		"file_type":  uploadedFile.FileType,
		"extension":  uploadedFile.Extension,
		"category":   uploadedFile.Category,
		"url":        "/" + relativePath,
		"created_at": uploadedFile.CreatedAt,
	})
}

func (h *FileHandler) Download(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	file, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(file.FilePath)
}

func (h *FileHandler) Delete(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
