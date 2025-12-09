package common

import "errors"

var (
	// 认证相关错误
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrUnauthorized       = errors.New("unauthorized access")

	// 用户相关错误
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user id")

	// 文件相关错误
	ErrFileNotFound        = errors.New("file not found")
	ErrFileToLarge         = errors.New("file size exceeds maximum limit")
	ErrInvalidFileType     = errors.New("invalid file type")
	ErrFileUploadFailed    = errors.New("file upload failed")
	ErrFileDeleteFailed    = errors.New("file delete failed")
	ErrInvalidFileName     = errors.New("invalid file name")
	ErrFilePathNotSafe     = errors.New("file path contains unsafe characters")

	// 资源相关错误
	ErrResourceNotFound  = errors.New("resource not found")
	ErrResourceForbidden = errors.New("access to resource is forbidden")
	ErrInvalidID         = errors.New("invalid id parameter")

	// 数据验证错误
	ErrInvalidInput      = errors.New("invalid input data")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrValidationFailed  = errors.New("validation failed")

	// 数据库相关错误
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")
	ErrRecordNotFound     = errors.New("record not found")
	ErrDuplicateEntry     = errors.New("duplicate entry")

	// 系统错误
	ErrInternalServer   = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrConfigurationError = errors.New("configuration error")
)
