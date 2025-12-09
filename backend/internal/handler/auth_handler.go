package handler

import (
	"net/http"
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token    string      `json:"token"`
    User     UserProfile `json:"user"`
    Message  string      `json:"message"`
}

type UserProfile struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
    // 认证功能已取消：返回默认用户信息与提示
    defaultUserID := middleware.GetCurrentUserID(c)
    var user model.User
    if err := database.DB.First(&user, defaultUserID).Error; err != nil {
        c.JSON(http.StatusOK, AuthResponse{Token: "", User: UserProfile{ID: 0, Username: "guest", Email: ""}, Message: "auth disabled"})
        return
    }

    c.JSON(http.StatusOK, AuthResponse{
        Token:   "",
        User:    UserProfile{ID: user.ID, Username: user.Username, Email: user.Email},
        Message: "auth disabled",
    })
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
    // 认证功能已取消：返回默认用户信息与提示
    defaultUserID := middleware.GetCurrentUserID(c)
    var user model.User
    if err := database.DB.First(&user, defaultUserID).Error; err != nil {
        c.JSON(http.StatusOK, AuthResponse{Token: "", User: UserProfile{ID: 0, Username: "guest", Email: ""}, Message: "auth disabled"})
        return
    }

    c.JSON(http.StatusOK, AuthResponse{
        Token:   "",
        User:    UserProfile{ID: user.ID, Username: user.Username, Email: user.Email},
        Message: "auth disabled",
    })
}

// GetProfile 获取当前用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
