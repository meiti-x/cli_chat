package repository

import "github.com/meiti-x/snapp_task/internal/models"

type MessageRepository interface {
	CreateMessage(message *models.Message) error
	GetUserMessage(message *models.Message, filter models.MessageFilter) (*models.Message, error)
}
