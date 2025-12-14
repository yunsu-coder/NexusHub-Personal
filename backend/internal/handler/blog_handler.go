package handler

import (
	"net/http"
	"nexushub-personal/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BlogHandler struct {
	DB *gorm.DB
}

func NewBlogHandler(db *gorm.DB) *BlogHandler {
	return &BlogHandler{DB: db}
}

// GetAllPosts returns all blog posts
func (h *BlogHandler) GetAllPosts(c *gin.Context) {
	var posts []model.Post
	userID := c.GetUint("user_id")

	if err := h.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// CreatePost creates a new blog post
func (h *BlogHandler) CreatePost(c *gin.Context) {
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.UserID = c.GetUint("user_id")
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// 自动生成摘要
	if post.Excerpt == "" {
		runes := []rune(post.Content)
		if len(runes) > 100 {
			post.Excerpt = string(runes[:100]) + "..."
		} else {
			post.Excerpt = post.Content
		}
	}
	// 默认封面
	if post.Cover == "" {
		post.Cover = "https://picsum.photos/seed/" + strconv.FormatInt(time.Now().Unix(), 10) + "/300/200"
	}
	// 默认作者
	if post.Author == "" {
		post.Author = "Admin" // 简化处理，实际可取当前用户名
	}

	if err := h.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost updates a blog post
func (h *BlogHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var post model.Post
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var req model.Post
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Tags = req.Tags
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	} else {
		// 重新生成摘要
		runes := []rune(post.Content)
		if len(runes) > 100 {
			post.Excerpt = string(runes[:100]) + "..."
		} else {
			post.Excerpt = post.Content
		}
	}
	if req.Cover != "" {
		post.Cover = req.Cover
	}
	post.UpdatedAt = time.Now()

	if err := h.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost deletes a blog post
func (h *BlogHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Post{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
