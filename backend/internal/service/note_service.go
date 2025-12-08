package service

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/model"
)

type NoteService struct{}

func NewNoteService() *NoteService {
	return &NoteService{}
}

func (s *NoteService) GetAll(userID uint) ([]model.Note, error) {
	var notes []model.Note
	err := database.DB.Where("user_id = ?", userID).Order("is_pinned DESC, updated_at DESC").Find(&notes).Error
	return notes, err
}

func (s *NoteService) GetByID(id, userID uint) (*model.Note, error) {
	var note model.Note
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error
	return &note, err
}

func (s *NoteService) Create(note *model.Note) error {
	return database.DB.Create(note).Error
}

func (s *NoteService) Update(note *model.Note) error {
	return database.DB.Save(note).Error
}

func (s *NoteService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Note{}).Error
}
