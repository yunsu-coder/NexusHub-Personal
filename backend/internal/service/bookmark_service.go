package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type BookmarkService struct{}

func NewBookmarkService() *BookmarkService {
	return &BookmarkService{}
}

func (s *BookmarkService) GetAll(userID uint) ([]model.Bookmark, error) {
	var bookmarks []model.Bookmark
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&bookmarks).Error
	return bookmarks, err
}

func (s *BookmarkService) GetByID(id, userID uint) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&bookmark).Error
	return &bookmark, err
}

func (s *BookmarkService) Create(bookmark *model.Bookmark) error {
	return database.DB.Create(bookmark).Error
}

func (s *BookmarkService) Update(bookmark *model.Bookmark) error {
	return database.DB.Save(bookmark).Error
}

func (s *BookmarkService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bookmark{}).Error
}
