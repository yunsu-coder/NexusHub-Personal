package constants

// 任务状态
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
)

// 任务优先级
const (
	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
)

// 文件类型
const (
	FileTypeFolder   = "folder"
	FileTypeImage    = "image"
	FileTypeDocument = "document"
	FileTypeCode     = "code"
	FileTypeArchive  = "archive"
	FileTypeVideo    = "video"
	FileTypeAudio    = "audio"
	FileTypeUnknown  = "unknown"
)

// 文件扩展名分类
var FileExtensions = map[string][]string{
	FileTypeImage:    {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico"},
	FileTypeDocument: {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md"},
	FileTypeCode:     {".js", ".ts", ".vue", ".jsx", ".tsx", ".go", ".py", ".java", ".c", ".cpp", ".html", ".css", ".json", ".xml"},
	FileTypeArchive:  {".zip", ".rar", ".7z", ".tar", ".gz"},
	FileTypeVideo:    {".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv"},
	FileTypeAudio:    {".mp3", ".wav", ".flac", ".aac", ".ogg"},
}

// 默认配置
const (
	DefaultPageSize = 20
	MaxFileSize     = 10 * 1024 * 1024 // 10MB
	TokenExpireHours = 24
)
