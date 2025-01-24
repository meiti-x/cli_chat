package storage

import (
	"fmt"
	"github.com/meiti-x/snapp_task/internal/models"
	"github.com/meiti-x/snapp_task/internal/repository"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func (r *messageRepository) CreateMessage(message *models.Message) error {
	if err := r.db.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

// TODO
func (r *messageRepository) GetUserMessage(message *models.Message, filter models.MessageFilter) (*models.Message, error) {
	return nil, nil
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{db}
}
