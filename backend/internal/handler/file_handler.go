package handler

import (
	"nexushub-personal/internal/common"
	"nexushub-personal/internal/logger"
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
		logger.Error("Failed to get all files: %v, user_id=%d", err, userID)
		common.InternalServerError(c, "Failed to retrieve files")
		return
	}

	logger.Info("Retrieved %d files for user_id=%d", len(files), userID)
	common.Success(c, files)
}

func (h *FileHandler) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid file ID parameter: %v", err)
		common.BadRequest(c, "Invalid file ID")
		return
	}

	file, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		logger.Warn("File not found: id=%d, user_id=%d", id, userID)
		common.NotFound(c, "File not found")
		return
	}

	common.Success(c, file)
}

func (h *FileHandler) GetByCategory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	category := c.Param("category")

	if category == "" {
		logger.Warn("Empty category parameter")
		common.BadRequest(c, "Category parameter is required")
		return
	}

	files, err := h.service.GetByCategory(category, userID)
	if err != nil {
		logger.Error("Failed to get files by category: %v, category=%s, user_id=%d", err, category, userID)
		common.InternalServerError(c, "Failed to retrieve files")
		return
	}

	logger.Info("Retrieved %d files for category=%s, user_id=%d", len(files), category, userID)
	common.Success(c, files)
}

func (h *FileHandler) Upload(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn("No file in upload request: %v, user_id=%d", err, userID)
		common.BadRequest(c, "No file uploaded")
		return
	}

	uploadedFile, err := h.service.Upload(file, userID)
	if err != nil {
		// 错误已在service层记录
		if err == common.ErrFileToLarge {
			common.BadRequest(c, err.Error())
		} else if err == common.ErrInvalidFileType || err == common.ErrInvalidFileName {
			common.BadRequest(c, err.Error())
		} else {
			common.InternalServerError(c, "File upload failed")
		}
		return
	}

	// Generate accessible URL path for the uploaded file
	relativePath := ""
	if idx := strings.Index(uploadedFile.FilePath, "uploads"); idx != -1 {
		relativePath = uploadedFile.FilePath[idx:]
		// Convert backslashes to forward slashes for URL
		relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	}

	response := gin.H{
		"id":         uploadedFile.ID,
		"file_name":  uploadedFile.FileName,
		"file_size":  uploadedFile.FileSize,
		"file_type":  uploadedFile.FileType,
		"extension":  uploadedFile.Extension,
		"category":   uploadedFile.Category,
		"url":        "/" + relativePath,
		"created_at": uploadedFile.CreatedAt,
	}

	common.Created(c, response)
}

func (h *FileHandler) Download(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid file ID for download: %v", err)
		common.BadRequest(c, "Invalid file ID")
		return
	}

	file, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		logger.Warn("File not found for download: id=%d, user_id=%d", id, userID)
		common.NotFound(c, "File not found")
		return
	}

	logger.Info("File download: id=%d, filename=%s, user_id=%d", id, file.FileName, userID)
	c.File(file.FilePath)
}

func (h *FileHandler) Delete(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid file ID for deletion: %v", err)
		common.BadRequest(c, "Invalid file ID")
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		// 错误已在service层记录
		if err == common.ErrFileNotFound {
			common.NotFound(c, "File not found")
		} else {
			common.InternalServerError(c, "Failed to delete file")
		}
		return
	}

	common.SuccessWithMessage(c, "File deleted successfully", nil)
}
