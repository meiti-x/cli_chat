package repository

import "github.com/meiti-x/snapp_task/internal/models"

type Message interface {
	CreateMessage(message *models.Message) (*models.Message, error)
	GetUserMessage(message *models.Message, filter models.MessageFilter) (*models.Message, error)
}
