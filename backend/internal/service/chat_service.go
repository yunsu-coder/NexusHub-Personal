package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type ChatService struct{}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) GetHistory(userID uint, limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	query := database.DB.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&messages).Error
	return messages, err
}

func (s *ChatService) Create(message *model.ChatMessage) error {
	return database.DB.Create(message).Error
}

func (s *ChatService) DeleteHistory(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.ChatMessage{}).Error
}
