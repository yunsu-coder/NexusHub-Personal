package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type CodeService struct{}

func NewCodeService() *CodeService {
	return &CodeService{}
}

func (s *CodeService) GetAll(userID uint) ([]model.CodeSnippet, error) {
	var snippets []model.CodeSnippet
	err := database.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&snippets).Error
	return snippets, err
}

func (s *CodeService) GetByID(id, userID uint) (*model.CodeSnippet, error) {
	var snippet model.CodeSnippet
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&snippet).Error
	return &snippet, err
}

func (s *CodeService) GetByLanguage(language string, userID uint) ([]model.CodeSnippet, error) {
	var snippets []model.CodeSnippet
	err := database.DB.Where("user_id = ? AND language = ?", userID, language).Order("updated_at DESC").Find(&snippets).Error
	return snippets, err
}

func (s *CodeService) Create(snippet *model.CodeSnippet) error {
	return database.DB.Create(snippet).Error
}

func (s *CodeService) Update(snippet *model.CodeSnippet) error {
	return database.DB.Save(snippet).Error
}

func (s *CodeService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.CodeSnippet{}).Error
}
