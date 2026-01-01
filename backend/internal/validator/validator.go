package validator

import (
	"fmt"
	"mime/multipart"
	"nexushub-personal/internal/common"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateFileUpload 验证文件上传
func ValidateFileUpload(fileHeader *multipart.FileHeader, maxSize int64, allowedExts []string) error {
	if fileHeader == nil {
		return common.ErrMissingRequiredField
	}

	// 检查文件大小
	if fileHeader.Size > maxSize {
		return fmt.Errorf("%w: maximum allowed size is %d bytes", common.ErrFileToLarge, maxSize)
	}

	if fileHeader.Size == 0 {
		return fmt.Errorf("%w: file is empty", common.ErrInvalidInput)
	}

	// 检查文件名
	if err := ValidateFileName(fileHeader.Filename); err != nil {
		return err
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	
	// 如果没有指定允许的文件类型，使用默认支持的文件类型列表
	if len(allowedExts) == 0 {
		defaultTypes := getDefaultAcceptedTypes()
		allowedExts = defaultTypes
	}
	
	if len(allowedExts) > 0 {
		allowed := false
		for _, allowedExt := range allowedExts {
			if ext == strings.ToLower(allowedExt) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("%w: %s is not allowed", common.ErrInvalidFileType, ext)
		}
	}

	return nil
}

// getDefaultAcceptedTypes 获取默认支持的文件类型
func getDefaultAcceptedTypes() []string {
	return []string{
		// 图片格式
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".ico", ".heic", ".tiff",
		// 文档格式
		".pdf", ".doc", ".docx", ".txt", ".md", ".rtf", ".xls", ".xlsx", ".ppt", ".pptx",
		// 代码格式
		".js", ".ts", ".vue", ".jsx", ".tsx", ".go", ".py", ".java", ".c", ".cpp", ".cs", ".php", ".rb", ".rs", ".sh", ".json", ".xml", ".yaml", ".yml",
		// 音频格式
		".mp3", ".wav", ".flac", ".ogg", ".aac", ".m4a",
		// 视频格式
		".mp4", ".webm", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".m4v",
		// 压缩格式
		".zip", ".rar", ".7z", ".tar", ".gz", ".bz2",
	}
}

// ValidateFileName 验证文件名安全性
func ValidateFileName(filename string) error {
	if filename == "" {
		return fmt.Errorf("%w: filename cannot be empty", common.ErrInvalidFileName)
	}

	// 检查文件名长度
	if len(filename) > 255 {
		return fmt.Errorf("%w: filename too long (max 255 characters)", common.ErrInvalidFileName)
	}

	// 检查危险字符
	dangerousChars := []string{"../", "..\\", "<", ">", ":", "\"", "|", "?", "*", "\x00"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("%w: filename contains unsafe character: %s", common.ErrFilePathNotSafe, char)
		}
	}

	// 检查是否只包含空格
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("%w: filename cannot be only whitespace", common.ErrInvalidFileName)
	}

	return nil
}

// ValidateID 验证ID参数
func ValidateID(id uint) error {
	if id == 0 {
		return common.ErrInvalidID
	}
	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return nil // 邮箱可以为空
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%w: invalid email format", common.ErrInvalidInput)
	}

	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("%w: username is required", common.ErrMissingRequiredField)
	}

	if len(username) < 3 {
		return fmt.Errorf("%w: username must be at least 3 characters", common.ErrInvalidInput)
	}

	if len(username) > 50 {
		return fmt.Errorf("%w: username must be less than 50 characters", common.ErrInvalidInput)
	}

	// 只允许字母、数字、下划线和连字符
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("%w: username can only contain letters, numbers, underscores and hyphens", common.ErrInvalidInput)
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("%w: password is required", common.ErrMissingRequiredField)
	}

	if len(password) < 6 {
		return fmt.Errorf("%w: password must be at least 6 characters", common.ErrInvalidInput)
	}

	if len(password) > 100 {
		return fmt.Errorf("%w: password must be less than 100 characters", common.ErrInvalidInput)
	}

	return nil
}

// ValidateStringLength 验证字符串长度
func ValidateStringLength(value, fieldName string, min, max int) error {
	length := len(value)

	if length < min {
		return fmt.Errorf("%w: %s must be at least %d characters", common.ErrInvalidInput, fieldName, min)
	}

	if length > max {
		return fmt.Errorf("%w: %s must be less than %d characters", common.ErrInvalidInput, fieldName, max)
	}

	return nil
}

// ValidateRequired 验证必填字段
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: %s is required", common.ErrMissingRequiredField, fieldName)
	}
	return nil
}

// SanitizeFilePath 清理文件路径
func SanitizeFilePath(path string) string {
	// 移除危险字符
	path = strings.ReplaceAll(path, "..", "")
	path = strings.ReplaceAll(path, "~", "")
	path = filepath.Clean(path)
	return path
}
