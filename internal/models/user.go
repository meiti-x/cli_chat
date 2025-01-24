package models

// Username is a custom type for user's username
type Username string

// User is a model for user
type User struct {
	ID       uint     `gorm:"primaryKey"`
	Username Username `gorm:"unique;not null"`
	Password string   `gorm:"not null"`
}
