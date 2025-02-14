package storage

import (
	"fmt"
	"github.com/meiti-x/cli_chat/internal/models"
	"github.com/meiti-x/cli_chat/internal/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	if r.db == nil {
		fmt.Println("Database connection is nil")
	}
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUsername(username *models.Username) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
