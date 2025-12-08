package model

import (
	"time"
	"gorm.io/gorm"
)

// User represents the single user of the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:100;not null;unique" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 密码不返回给前端
	Email     string         `gorm:"size:255;unique" json:"email"`
	Nickname  string         `gorm:"size:100" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Bio       string         `gorm:"size:500" json:"bio"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Note represents a note/memo
type Note struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:longtext" json:"content"`
	Tags      string         `gorm:"size:500" json:"tags"`
	IsPinned  bool           `gorm:"default:false" json:"is_pinned"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// File represents uploaded files
type File struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	FileName    string         `gorm:"size:255;not null" json:"file_name"`
	FilePath    string         `gorm:"size:500;not null" json:"file_path"`
	FileSize    int64          `gorm:"not null" json:"file_size"`
	FileType    string         `gorm:"size:100" json:"file_type"` // video, image, document, code, etc
	MimeType    string         `gorm:"size:100" json:"mime_type"`
	Extension   string         `gorm:"size:20" json:"extension"`
	Thumbnail   string         `gorm:"size:500" json:"thumbnail"`
	Category    string         `gorm:"size:50" json:"category"` // media, document, code, archive, etc
	Description string         `gorm:"size:500" json:"description"`
	Tags        string         `gorm:"size:500" json:"tags"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CodeSnippet represents code snippets with syntax highlighting
type CodeSnippet struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Language  string         `gorm:"size:50;not null" json:"language"` // go, cpp, python, markdown, etc
	Code      string         `gorm:"type:longtext;not null" json:"code"`
	Tags      string         `gorm:"size:500" json:"tags"`
	IsFavorite bool          `gorm:"default:false" json:"is_favorite"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Theme represents user theme preferences
type Theme struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	UserID            uint      `gorm:"not null;unique;index" json:"user_id"`
	ThemeName         string    `gorm:"size:50;default:'dark'" json:"theme_name"` // dark, light, custom
	PrimaryColor      string    `gorm:"size:20;default:'#000000'" json:"primary_color"`
	SecondaryColor    string    `gorm:"size:20;default:'#ffffff'" json:"secondary_color"`
	BackgroundImage   string    `gorm:"size:500" json:"background_image"`
	BackgroundOpacity float32   `gorm:"default:1.0" json:"background_opacity"` // 0.0 to 1.0, default fully opaque
	BackgroundMusic   string    `gorm:"size:500" json:"background_music"`
	MusicVolume       float32   `gorm:"default:0.5" json:"music_volume"`
	CustomCSS         string    `gorm:"type:text" json:"custom_css"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ChatMessage represents AI chat messages
type ChatMessage struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Role      string         `gorm:"size:20;not null" json:"role"` // user, assistant
	Content   string         `gorm:"type:longtext;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Task represents todo tasks
type Task struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"size:20;default:'pending'" json:"status"`  // pending, in_progress, completed
	Priority    string         `gorm:"size:20;default:'medium'" json:"priority"` // low, medium, high
	Category    string         `gorm:"size:100" json:"category"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Bookmark represents saved bookmarks
type Bookmark struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	URL         string         `gorm:"size:1000;not null" json:"url"`
	Description string         `gorm:"size:500" json:"description"`
	Tags        string         `gorm:"size:500" json:"tags"`
	Favicon     string         `gorm:"size:500" json:"favicon"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
