package db

import (
	"fmt"
	"github.com/meiti-x/snapp_task/config"
	"github.com/meiti-x/snapp_task/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// InitDB initializes the PostgreSQL database and GORM.
func InitDB(conf *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.Database.Host, conf.Database.User, conf.Database.Pass, conf.Database.Name, conf.Database.Port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Auto-migrate the User and Message models
	if err := database.AutoMigrate(&models.User{}, &models.Message{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, err
	}
	return database, nil

}
