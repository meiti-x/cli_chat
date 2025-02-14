package storage

import (
	"github.com/meiti-x/cli_chat/internal/models"
	"github.com/meiti-x/cli_chat/internal/repository"
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

// TODO: Implement the rest of the methods
func (r *messageRepository) GetUserMessage(message *models.Message, filter models.MessageFilter) (*models.Message, error) {
	return nil, nil
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{db}
}
