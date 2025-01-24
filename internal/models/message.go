package models

import "time"

// TODO: separate db models from domain models

// Message is a model for message
type Message struct {
	ID        uint      `gorm:"primaryKey"`
	Username  Username  `gorm:"not null"`
	Chatroom  string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// MessageFilter is a filter for message
type MessageFilter struct {
	Chatroom string `gorm:"not null"`
}
