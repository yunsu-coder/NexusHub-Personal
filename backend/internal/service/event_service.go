package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

// GetAll 获取用户所有事件，支持按日期范围过滤
func (s *EventService) GetAll(userID uint, startDate, endDate string) ([]model.Event, error) {
	var events []model.Event
	query := database.DB.Where("user_id = ?", userID)
	
	// 按日期范围过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		query = query.Where("date >= ?", startDate)
	} else if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	
	err := query.Order("date ASC, start_time ASC").Find(&events).Error
	return events, err
}

// GetByID 根据ID获取事件
func (s *EventService) GetByID(id, userID uint) (*model.Event, error) {
	var event model.Event
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&event).Error
	return &event, err
}

// Create 创建新事件
func (s *EventService) Create(event *model.Event) error {
	return database.DB.Create(event).Error
}

// Update 更新事件
func (s *EventService) Update(event *model.Event) error {
	return database.DB.Save(event).Error
}

// Delete 删除事件
func (s *EventService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Event{}).Error
}
