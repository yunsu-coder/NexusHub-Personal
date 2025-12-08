package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type TaskService struct{}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) GetAll(userID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := database.DB.Where("user_id = ?", userID).Order("priority DESC, created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) GetByID(id, userID uint) (*model.Task, error) {
	var task model.Task
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	return &task, err
}

func (s *TaskService) Create(task *model.Task) error {
	return database.DB.Create(task).Error
}

func (s *TaskService) Update(task *model.Task) error {
	return database.DB.Save(task).Error
}

func (s *TaskService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Task{}).Error
}
