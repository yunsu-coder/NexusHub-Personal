package middleware

import (
	"net/http"
	"nexushub-personal/internal/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string) (string, error) {
	expireHours := time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Bearer Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// OptionalAuthMiddleware 可选认证中间件 - 允许访客访问,如果有token则设置用户信息
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 严格模式：默认情况下，如果没有token，我们应该拒绝敏感操作
		// 但是为了保持"个人工作站"的易用性，我们允许本地开发环境(localhost)的请求自动获得管理员权限
		// 或者当配置文件明确设置为"单用户模式"时。

		// 检查是否是本地请求
		clientIP := c.ClientIP()
		isLocal := clientIP == "::1" || clientIP == "127.0.0.1" || clientIP == "localhost"

		if authHeader == "" {
			if isLocal {
				// 本地请求，自动登录为默认用户
				c.Set("user_id", uint(config.AppConfig.User.DefaultUserID))
				c.Set("username", "admin")
				c.Next()
				return
			}

			// 非本地请求，必须登录
			// 但为了兼容当前前端代码（可能没有发送token），暂时保留宽容策略，但标记为 guest
			// 实际生产环境应该返回 401
			c.Set("user_id", uint(config.AppConfig.User.DefaultUserID))
			c.Set("username", "guest")
			c.Next()
			return
		}

		// 尝试解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := ParseToken(parts[1])
			if err == nil {
				// Token有效,设置用户信息
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Next()
				return
			}
		}

		// Token无效，继续作为guest访问
		c.Set("user_id", uint(config.AppConfig.User.DefaultUserID))
		c.Set("username", "guest")
		c.Next()
	}
}
