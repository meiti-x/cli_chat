package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheMiss             = errors.New("cache miss")
	ErrRedisClose            = errors.New("failed to close Redis")
	ErrRedisInit             = errors.New("failed to initialize Redis connection")
	ErrRedisConnectionFailed = errors.New("failed to connect to Redis")
	ErrRedisOperationFailed  = errors.New("redis operation failed")
)

type Provider interface {
	// Set stores data in the cache with a given key and time-to-live (TTL).
	Set(ctx context.Context, key string, ttl time.Duration, data []byte) error

	// Get fetches data from the cache by its key.
	Get(ctx context.Context, key string) ([]byte, error)

	// Del deletes a key and its associated data from the cache.
	Del(ctx context.Context, key string) error

	// GetSetMembers fetches all members of a set stored at the given key.
	GetSetMembers(ctx context.Context, key string) ([]string, error)

	// RemoveSetMember removes a specific member from a set stored at the given key.
	RemoveSetMember(ctx context.Context, key string, value string) error

	// GetSetSize returns the number of elements in a set stored at the given key.
	GetSetSize(ctx context.Context, key string) (int64, error)

	// AddSetMember adds a new member to a set stored at the given key.
	AddSetMember(ctx context.Context, key string, value string) error

	// CloseConnection gracefully closes the Redis connection.
	CloseConnection() error
}
