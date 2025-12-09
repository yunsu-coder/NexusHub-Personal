package utils

import (
	"nexushub-personal/internal/constants"
	"path/filepath"
	"strings"
)

// GetFileType 根据文件扩展名获取文件类型
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	for fileType, extensions := range constants.FileExtensions {
		for _, e := range extensions {
			if e == ext {
				return fileType
			}
		}
	}

	return constants.FileTypeUnknown
}

// IsImageFile 判断是否为图片文件
func IsImageFile(filename string) bool {
	return GetFileType(filename) == constants.FileTypeImage
}

// IsCodeFile 判断是否为代码文件
func IsCodeFile(filename string) bool {
	return GetFileType(filename) == constants.FileTypeCode
}

// IsDocumentFile 判断是否为文档文件
func IsDocumentFile(filename string) bool {
	return GetFileType(filename) == constants.FileTypeDocument
}
