package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"nexushub-personal/internal/common"
	"nexushub-personal/internal/config"
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/logger"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/validator"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) GetAll(userID uint) ([]model.File, error) {
	var files []model.File
	err := database.DB.Where("user_id IN (?, ?)", userID, 0).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (s *FileService) GetByID(id, userID uint) (*model.File, error) {
	var file model.File
	err := database.DB.Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).First(&file).Error
	return &file, err
}

func (s *FileService) GetByCategory(category string, userID uint) ([]model.File, error) {
	var files []model.File
	err := database.DB.Where("category = ? AND user_id IN (?, ?)", category, userID, 0).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (s *FileService) Upload(fileHeader *multipart.FileHeader, userID uint) (*model.File, error) {
	// 允许访客用户上传文件（userID = 0）
	// 其他验证保持不变

	// 验证文件上传 - 不限制扩展名类型,只限制大小
	maxSize := config.AppConfig.Storage.MaxUploadSize
	if err := validator.ValidateFileUpload(fileHeader, maxSize, nil); err != nil {
		logger.Warn("File validation failed: %v, file: %s, size: %d", err, fileHeader.Filename, fileHeader.Size)
		return nil, err
	}

	// Get file extension and type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	category := s.getCategoryByExtension(ext)

	// 清理和验证存储路径
	baseDir := validator.SanitizeFilePath(config.AppConfig.Storage.Path)
	uploadDir := filepath.Join(baseDir, "uploads", category)

	// Create storage directory if not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("Failed to create upload directory: %v, path: %s", err, uploadDir)
		return nil, fmt.Errorf("%w: failed to create storage directory", common.ErrInternalServer)
	}

	// Generate unique filename to avoid conflicts
	timestamp := time.Now().Unix()
	safeFilename := validator.SanitizeFilePath(fileHeader.Filename)
	filename := fmt.Sprintf("%d_%s", timestamp, safeFilename)
	filePath := filepath.Join(uploadDir, filename)

	// 开始数据库事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to begin transaction: %v", tx.Error)
		return nil, fmt.Errorf("%w: database transaction error", common.ErrInternalServer)
	}

	// 确保事务正确处理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			os.Remove(filePath) // 清理已上传的文件
			logger.Error("Panic during file upload: %v", r)
		}
	}()

	// Open and validate uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to open uploaded file: %v, filename: %s", err, fileHeader.Filename)
		return nil, fmt.Errorf("%w: cannot open uploaded file", common.ErrFileUploadFailed)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to create destination file: %v, path: %s", err, filePath)
		return nil, fmt.Errorf("%w: cannot create destination file", common.ErrFileUploadFailed)
	}
	defer dst.Close()

	// Copy file content with size tracking
	written, err := io.Copy(dst, src)
	if err != nil {
		tx.Rollback()
		os.Remove(filePath) // 清理失败的文件
		logger.Error("Failed to copy file content: %v, filename: %s", err, fileHeader.Filename)
		return nil, fmt.Errorf("%w: file copy failed", common.ErrFileUploadFailed)
	}

	// 验证写入的字节数
	if written != fileHeader.Size {
		tx.Rollback()
		os.Remove(filePath)
		logger.Error("File size mismatch: expected %d, got %d, filename: %s", fileHeader.Size, written, fileHeader.Filename)
		return nil, fmt.Errorf("%w: file size mismatch", common.ErrFileUploadFailed)
	}

	// Create file record in database
	file := &model.File{
		UserID:    userID,
		FileName:  fileHeader.Filename,
		FilePath:  filePath,
		FileSize:  fileHeader.Size,
		FileType:  fileHeader.Header.Get("Content-Type"),
		MimeType:  fileHeader.Header.Get("Content-Type"),
		Extension: ext,
		Category:  category,
	}

	if err := tx.Create(file).Error; err != nil {
		tx.Rollback()
		os.Remove(filePath) // 数据库插入失败时清理文件
		logger.Error("Failed to create file record in database: %v, filename: %s", err, fileHeader.Filename)
		return nil, fmt.Errorf("%w: database insert failed", common.ErrInternalServer)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		os.Remove(filePath) // 提交失败时清理文件
		logger.Error("Failed to commit transaction: %v, filename: %s", err, fileHeader.Filename)
		return nil, fmt.Errorf("%w: transaction commit failed", common.ErrInternalServer)
	}

	logger.Info("File uploaded successfully: id=%d, filename=%s, size=%d, category=%s, user_id=%d",
		file.ID, file.FileName, file.FileSize, file.Category, userID)

	return file, nil
}

func (s *FileService) Rename(id, userID uint, newName string) error {
	// 验证文件名
	safeName := validator.SanitizeFilePath(newName)
	if safeName == "" {
		return common.ErrInvalidFileName
	}

	var file model.File
	if err := database.DB.Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).First(&file).Error; err != nil {
		return common.ErrFileNotFound
	}

	// 保持扩展名
	ext := filepath.Ext(file.FileName)
	if filepath.Ext(safeName) == "" {
		safeName += ext
	}

	oldPath := file.FilePath
	newPath := filepath.Join(filepath.Dir(oldPath), safeName)

	if err := os.Rename(oldPath, newPath); err != nil {
		logger.Error("Failed to rename physical file: %v", err)
		return common.ErrInternalServer
	}

	file.FileName = safeName
	file.FilePath = newPath
	return database.DB.Save(&file).Error
}

func (s *FileService) Delete(id, userID uint) error {
	// 验证ID
	if err := validator.ValidateID(id); err != nil {
		logger.Error("Invalid file ID for deletion: %v", err)
		return err
	}

	if err := validator.ValidateID(userID); err != nil {
		logger.Error("Invalid user ID for file deletion: %v", err)
		return err
	}

	// 开始事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to begin transaction for file deletion: %v", tx.Error)
		return fmt.Errorf("%w: database transaction error", common.ErrInternalServer)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Error("Panic during file deletion: %v", r)
		}
	}()

	// 查询文件记录
	var file model.File
	if err := tx.Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).First(&file).Error; err != nil {
		tx.Rollback()
		logger.Warn("File not found for deletion: id=%d, user_id=%d", id, userID)
		return common.ErrFileNotFound
	}

	// 删除数据库记录(软删除)
	result := tx.Delete(&file)
	if result.Error != nil {
		tx.Rollback()
		logger.Error("Failed to delete file record from database: %v, id=%d", result.Error, id)
		return fmt.Errorf("%w: database delete failed", common.ErrFileDeleteFailed)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		logger.Warn("No file record deleted (likely permission issue): id=%d, user_id=%d", id, userID)
		return common.ErrFileNotFound
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit file deletion transaction: %v, id=%d", err, id)
		return fmt.Errorf("%w: transaction commit failed", common.ErrInternalServer)
	}

	// 删除物理文件(在事务提交后,即使失败也不影响数据库)
	if err := os.Remove(file.FilePath); err != nil {
		logger.Warn("Failed to delete physical file (record deleted): %v, path=%s", err, file.FilePath)
		// 不返回错误,因为数据库记录已经删除
	}

	// 删除缩略图(如果存在)
	if file.Thumbnail != "" {
		if err := os.Remove(file.Thumbnail); err != nil {
			logger.Warn("Failed to delete thumbnail (ignored): %v, path=%s", err, file.Thumbnail)
		}
	}

	logger.Info("File deleted successfully: id=%d, filename=%s, user_id=%d", id, file.FileName, userID)
	return nil
}

func (s *FileService) getCategoryByExtension(ext string) string {
	mediaExts := map[string]bool{
		".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".webm": true,
		".mp3": true, ".wav": true, ".flac": true, ".aac": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true, ".svg": true,
	}

	documentExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true, ".txt": true, ".rtf": true,
	}

	codeExts := map[string]bool{
		".go": true, ".cpp": true, ".c": true, ".h": true, ".py": true,
		".js": true, ".ts": true, ".java": true, ".rs": true, ".php": true,
		".html": true, ".css": true, ".json": true, ".xml": true, ".yaml": true,
		".md": true, ".sh": true, ".sql": true,
	}

	archiveExts := map[string]bool{
		".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
	}

	if mediaExts[ext] {
		return "media"
	} else if documentExts[ext] {
		return "document"
	} else if codeExts[ext] {
		return "code"
	} else if archiveExts[ext] {
		return "archive"
	}

	return "other"
}
