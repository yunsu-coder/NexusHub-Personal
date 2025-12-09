package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type CollectionService struct{}

func NewCollectionService() *CollectionService {
	return &CollectionService{}
}

func (s *CollectionService) GetAll(userID uint, collectionType string) ([]model.Collection, error) {
	var collections []model.Collection
	query := database.DB.Where("user_id = ?", userID)
	
	// 按类型过滤
	if collectionType != "" {
		query = query.Where("type = ?", collectionType)
	}
	
	err := query.Order("created_at DESC").Find(&collections).Error
	return collections, err
}

func (s *CollectionService) GetByID(id, userID uint) (*model.Collection, error) {
	var collection model.Collection
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&collection).Error
	return &collection, err
}

func (s *CollectionService) Create(collection *model.Collection) error {
	return database.DB.Create(collection).Error
}

func (s *CollectionService) Update(collection *model.Collection) error {
	return database.DB.Save(collection).Error
}

func (s *CollectionService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Collection{}).Error
}
