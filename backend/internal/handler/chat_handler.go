package handler

import (
	"net/http"
	"nexushub-personal/internal/middleware"
	"nexushub-personal/internal/model"
	"nexushub-personal/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service *service.ChatService
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		service: service.NewChatService(),
	}
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	messages, err := h.service.GetHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save user message
	userMessage := &model.ChatMessage{
		UserID:  userID,
		Role:    "user",
		Content: req.Content,
	}
	if err := h.service.Create(userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Integrate with AI API
	// For now, return a simple echo response
	aiResponse := &model.ChatMessage{
		UserID:  userID,
		Role:    "assistant",
		Content: "AI功能待接入，当前为回显模式: " + req.Content,
	}
	if err := h.service.Create(aiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_message": userMessage,
		"ai_response":  aiResponse,
	})
}

func (h *ChatHandler) ClearHistory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if err := h.service.DeleteHistory(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Chat history cleared successfully"})
}
