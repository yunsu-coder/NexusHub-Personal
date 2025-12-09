package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"

	"gorm.io/gorm"
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
	// 先尝试获取用户的现有主题
	var existingTheme model.Theme
	result := database.DB.Where("user_id = ?", theme.UserID).First(&existingTheme)

	if result.Error != nil {
		// 如果主题不存在且错误是记录未找到，则创建一个新的
		if result.Error == gorm.ErrRecordNotFound {
			return database.DB.Create(theme).Error
		}
		// 如果是其他错误，则返回
		return result.Error
	}

	// 合并主题属性
	existingTheme.ThemeName = theme.ThemeName
	existingTheme.PrimaryColor = theme.PrimaryColor
	existingTheme.SecondaryColor = theme.SecondaryColor
	existingTheme.ThemeTemplate = theme.ThemeTemplate
	existingTheme.BackgroundMusic = theme.BackgroundMusic
	existingTheme.MusicVolume = theme.MusicVolume
	existingTheme.CustomCSS = theme.CustomCSS

	// 保存更新后的主题
	return database.DB.Save(existingTheme).Error
}
