package repository

import "github.com/meiti-x/snapp_task/internal/models"

type UserRepository interface {
	GetUser(user *models.Username) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
}
