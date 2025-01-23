package models

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	Username  Username  `gorm:"not null"`
	Chatroom  string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MessageFilter struct {
	Chatroom string `gorm:"not null"`
}
