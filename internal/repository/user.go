package repository

import (
	"context"
	"github.com/meiti-x/cli_chat/internal/models"
)

type UserRepository interface {
	GetUserByUsername(user *models.Username) (*models.User, error)
	CreateUser(user *models.User) error
}

type UserRedisRepository interface {
	AddUserToChatroomSet(ctx context.Context, onlineUsersKey string, clientIP string) error
	RemoveUserFromChatroomSet(ctx context.Context, onlineUsersKey string, clientIP string) error
	GetTotalUserInChatroom(ctx context.Context, onlineUsersKey string, clientIP string) error
}
