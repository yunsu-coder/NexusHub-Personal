package service

import (
	"nexushub-personal/internal/common"
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type CollectionService struct{}

func NewCollectionService() *CollectionService {
	return &CollectionService{}
}

func (s *CollectionService) GetAll(userID uint, collectionType string) ([]model.Collection, error) {
	var collections []model.Collection
	query := database.DB.Where("user_id IN (?, ?)", userID, 0)

	// 按类型过滤
	if collectionType != "" {
		query = query.Where("type = ?", collectionType)
	}

	err := query.Order("created_at DESC").Find(&collections).Error
	return collections, err
}

func (s *CollectionService) GetByID(id, userID uint) (*model.Collection, error) {
	var collection model.Collection
	err := database.DB.Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).First(&collection).Error
	return &collection, err
}

func (s *CollectionService) Create(collection *model.Collection) error {
	return database.DB.Create(collection).Error
}

func (s *CollectionService) Update(collection *model.Collection) error {
	return database.DB.Save(collection).Error
}

func (s *CollectionService) Delete(id, userID uint) error {
	result := database.DB.Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).Delete(&model.Collection{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrFileNotFound
	}
	return nil
}
