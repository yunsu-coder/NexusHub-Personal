package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type ThemeService struct{}

func NewThemeService() *ThemeService {
	return &ThemeService{}
}

func (s *ThemeService) Get(userID uint) (*model.Theme, error) {
	var theme model.Theme
	err := database.DB.Where("user_id = ?", userID).First(&theme).Error
	return &theme, err
}

func (s *ThemeService) Update(theme *model.Theme) error {
	return database.DB.Save(theme).Error
}
