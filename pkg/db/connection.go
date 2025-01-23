package db

import (
	"github.com/meiti-x/snapp_task/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// initDB initializes the PostgreSQL database and GORM.
func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=chatroom_db port=5432 sslmode=disable"
	//var database *gorm.DB

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the User and Message models
	if err := database.AutoMigrate(&models.User{}, &models.Message{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, err
	}
	return database, nil

}
