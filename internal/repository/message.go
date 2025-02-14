package repository

import "github.com/meiti-x/cli_chat/internal/models"

// MessageRepository is an interface for message repository
type MessageRepository interface {
	CreateMessage(message *models.Message) error
	GetUserMessage(message *models.Message, filter models.MessageFilter) (*models.Message, error)
}
