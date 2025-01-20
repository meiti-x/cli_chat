package users

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/meiti-x/snapp_chal/socket/utils"
)

var (
	ErrInvalidUserIDCreate = errors.New("userID is required and must be greater than 0")
)

type GenderEnum string
type UserID uint

const (
	Male   GenderEnum = "male"
	Female GenderEnum = "female"
)

// User represents the database model for a users
type User struct {
	ID                     uint        `gorm:"primaryKey"`
	FirstName              string      `gorm:"size:255"`
	LastName               string      `gorm:"size:255"`
	Birthdate              string      `gorm:"size:255"`
	City                   string      `gorm:"size:255"`
	NationalCode           string      `gorm:"size:10;unique"`
	Gender                 *GenderEnum `gorm:"type:gender_enum"`
	Email                  string      `gorm:"unique;size:255"`
	Password               string      `gorm:"not null"`
	IsActive               bool        `gorm:"not null"`
	WalletBalance          int64
	MaxQuestionnairesCount int `gorm:"null"`
	CreatedAt              time.Time
	UpdatedAt              time.Time `gorm:"not null"`
}

// TwoFACode stores 2FA codes for users
type TwoFACode struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserDTO represents the data needed to create a new users
type CreateUserDTO struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	NationalCode string `json:"national_code" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

// UpdateUserDTO represents the data needed to update an existing users
type UpdateUserDTO struct {
	FirstName    *string     `json:"first_name" validate:"required"`
	LastName     *string     `json:"last_name" validate:"required"`
	Birthdate    *string     `json:"birthdate" validate:"required"`
	Gender       *GenderEnum `json:"gender" validate:"required,oneof=male female"`
	City         *string     `json:"city" validate:"required"`
	NationalCode *string     `json:"national_code" validate:"required"`
}

// IncreaseWalletBalanceDTO represents the data needed to update credit users
type IncreaseWalletBalanceDTO struct {
	Value string `json:"value" validate:"required"`
}

// LoginRequest represents users login data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication responses
type AuthResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TwoFACodeSent bool   `json:"two_fa_code_sent"`
}

// Verify2FACodeRequest validates 2FA code
type Verify2FACodeRequest struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

// UserResponse represents the users data returned in API responses
type UserResponse struct {
	ID            UserID `json:"id"`
	Name          string `json:"name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	NationalCode  string `json:"national_code"`
	Gender        string `json:"gender"`
	Birthdate     string `json:"birthdate"`
	City          string `json:"city"`
	WalletBalance int64  `json:"wallet_balance"`
}

// PublicUserResponse represents the users data returned in API responses
type PublicUserResponse struct {
	ID     UserID `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

// GetFullName returns the full name of a users
func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) GetGender() string {
	if u.Gender == nil {
		return "unknown"
	}
	return string(*u.Gender)
}

// ToUserResponse maps a User model to a UserResponse DTO
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID: UserID(user.ID),
		//Name:          users.GetFullName(),
		Name:         "Test",
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		NationalCode: user.NationalCode,
		//Gender:        users.GetGender(),
		Gender:        "Male",
		City:          user.City,
		Birthdate:     user.Birthdate,
		WalletBalance: user.WalletBalance,
	}
}

// ToPublicUserResponse maps a User model to a UserResponse DTO
func ToPublicUserResponse(user *User) *PublicUserResponse {
	return &PublicUserResponse{
		ID: UserID(user.ID),
		//Name:   users.GetFullName(),
		//Gender: users.GetGender(),
	}
}

// ToUserModel maps a CreateUserDTO to a User model
func ToUserModel(dto *CreateUserDTO) *User {
	return &User{
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		NationalCode: dto.NationalCode,
		Password:     dto.Password,
	}
}

// ToUserModelForUpdate maps a UpdateUserDTO to a User model
func ToUserModelForUpdate(user User, dto *UpdateUserDTO) User {
	if dto.FirstName != nil {
		user.FirstName = *dto.FirstName
	}

	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}

	if dto.Birthdate != nil {
		user.Birthdate = *dto.Birthdate
	}

	if dto.Gender != nil {
		user.Gender = dto.Gender
	}

	if dto.City != nil {
		user.City = *dto.City
	}

	if dto.NationalCode != nil {
		user.NationalCode = *dto.NationalCode
	}

	return user
}

// Validate checks the User struct for common validation rules.
func (u *User) Validate() error {
	if strings.TrimSpace(u.FirstName) == "" {
		return errors.New("first name is required")
	}
	if strings.TrimSpace(u.LastName) == "" {
		return errors.New("last name is required")
	}
	if !utils.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	isValidNationalCode, err := utils.IsValidNationalCode(u.NationalCode)
	if err != nil {
		return errors.New("national code validation failed")
	}
	if !isValidNationalCode {
		return errors.New("national code is not valid")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func (u *UpdateUserDTO) ValidateForUpdate() error {
	if u.FirstName != nil {
		if strings.TrimSpace(*u.FirstName) == "" {
			return errors.New("first name can not empty string")
		}
	}

	if u.LastName != nil {
		if strings.TrimSpace(*u.LastName) == "" {
			return errors.New("last name can not empty string")
		}
	}

	if u.NationalCode != nil {
		isValidNationalCode, err := utils.IsValidNationalCode(*u.NationalCode)
		if err != nil {
			return errors.New("national code validation failed")
		}

		if !isValidNationalCode {
			return errors.New("national code is not valid")
		}
	}

	if u.Birthdate != nil {
		isValid, message := utils.IsValidBirthdate(*u.Birthdate)
		if !isValid {
			return errors.New("birthdate - " + message)
		}
	}

	if u.City != nil {
		isValid, message := utils.IsValidCity(*u.City)
		if !isValid {
			return errors.New("city - " + message)
		}
	}

	if u.Gender != nil {
		if *u.Gender != Male && *u.Gender != Female {
			return errors.New("gender value is invalid")
		}
	}

	return nil
}
