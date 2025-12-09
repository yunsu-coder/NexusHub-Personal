package constants

// HTTP 状态码相关错误消息
const (
	ErrInvalidID       = "Invalid ID"
	ErrInvalidRequest  = "Invalid request"
	ErrUnauthorized    = "Unauthorized"
	ErrForbidden       = "Forbidden"
	ErrNotFound        = "Resource not found"
	ErrInternalServer  = "Internal server error"
	ErrDatabaseError   = "Database error"
)

// 资源类型错误消息
const (
	ErrNoteNotFound     = "Note not found"
	ErrTaskNotFound     = "Task not found"
	ErrBookmarkNotFound = "Bookmark not found"
	ErrCodeNotFound     = "Code snippet not found"
	ErrFileNotFound     = "File not found"
	ErrUserNotFound     = "User not found"
)

// 操作成功消息
const (
	MsgCreatedSuccess = "Created successfully"
	MsgUpdatedSuccess = "Updated successfully"
	MsgDeletedSuccess = "Deleted successfully"
	MsgLoginSuccess   = "Login successful"
	MsgRegisterSuccess = "Registration successful"
)

// 业务逻辑错误消息
const (
	ErrUsernameExists    = "Username already exists"
	ErrInvalidCredentials = "Invalid username or password"
	ErrTokenExpired      = "Token has expired"
	ErrTokenInvalid      = "Invalid token"
	ErrPermissionDenied  = "Permission denied"
)

// 文件操作错误消息
const (
	ErrFileUpload    = "File upload failed"
	ErrFileDownload  = "File download failed"
	ErrFileDelete    = "File delete failed"
	ErrFileTooLarge  = "File size exceeds limit"
	ErrInvalidFileType = "Invalid file type"
)
