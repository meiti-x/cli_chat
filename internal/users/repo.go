package users

import (
	"time"
)

// IUserRepository defines the repository methods
type IUserRepository interface {
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id UserID) error
	GetUserByEmail(email string) (*User, error)
	GetUserByNationalCode(nationalCode string) (*User, error)
	StoreTwoFACode(email string, code string, expiresAt time.Time) error
	GetTwoFACode(email string) (*TwoFACode, error)
	RemoveTwoFACode(email string) error
}
