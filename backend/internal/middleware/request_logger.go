package middleware

import (
	"time"

	"nexushub-personal/internal/logger"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		userID, _ := c.Get("user_id")

		logger.Info("%s %s | %d | %v | user_id=%v", method, path, status, duration, userID)
	}
}
