package models

type Username string
type User struct {
	ID       uint     `gorm:"primaryKey"`
	Username Username `gorm:"unique;not null"`
	Password string   `gorm:"not null"`
}
