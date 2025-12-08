package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"nexushub-personal/internal/config"
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
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
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (s *FileService) GetByID(id, userID uint) (*model.File, error) {
	var file model.File
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&file).Error
	return &file, err
}

func (s *FileService) GetByCategory(category string, userID uint) ([]model.File, error) {
	var files []model.File
	err := database.DB.Where("user_id = ? AND category = ?", userID, category).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (s *FileService) Upload(fileHeader *multipart.FileHeader, userID uint) (*model.File, error) {
	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Get file extension and type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	category := s.getCategoryByExtension(ext)

	// Create storage directory if not exists
	uploadDir := filepath.Join(config.AppConfig.Storage.Path, "uploads", category)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Create file record
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

	if err := database.DB.Create(file).Error; err != nil {
		// Remove uploaded file if database insert fails
		os.Remove(filePath)
		return nil, err
	}

	return file, nil
}

func (s *FileService) Delete(id, userID uint) error {
	var file model.File
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&file).Error; err != nil {
		return err
	}

	// Delete physical file
	if err := os.Remove(file.FilePath); err != nil {
		// Continue even if file deletion fails
	}

	// Delete thumbnail if exists
	if file.Thumbnail != "" {
		os.Remove(file.Thumbnail)
	}

	// Delete database record
	return database.DB.Delete(&file).Error
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
