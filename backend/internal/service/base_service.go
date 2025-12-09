package service

import (
	"gorm.io/gorm"
)

// BaseService 通用CRUD服务基类
type BaseService[T any] struct {
	db *gorm.DB
}

// NewBaseService 创建基础服务
func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{
		db: db,
	}
}

// GetAll 获取所有记录
func (s *BaseService[T]) GetAll(userID uint) ([]T, error) {
	var items []T
	err := s.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&items).Error
	return items, err
}

// GetByID 根据ID获取记录
func (s *BaseService[T]) GetByID(id, userID uint) (*T, error) {
	var item T
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	return &item, err
}

// Create 创建记录
func (s *BaseService[T]) Create(entity *T) error {
	return s.db.Create(entity).Error
}

// Update 更新记录
func (s *BaseService[T]) Update(entity *T) error {
	return s.db.Save(entity).Error
}

// Delete 删除记录
func (s *BaseService[T]) Delete(id, userID uint) error {
	var item T
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&item).Error
}
