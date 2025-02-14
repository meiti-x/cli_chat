package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/meiti-x/cli_chat/config"
	c "github.com/meiti-x/cli_chat/pkg/cache"
	"time"
)

// redisCacheClient is an implementation of the CacheProvider interface using Redis as the backend.
type redisCacheClient struct {
	client *redis.Client
}

func NewRedisCache(conf *config.Config) (c.Provider, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &redisCacheClient{
		client: rdb,
	}, nil
}

func (r *redisCacheClient) AddSetMember(ctx context.Context, key string, value string) error {
	err := r.client.SAdd(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCacheClient) RemoveSetMember(ctx context.Context, key string, value string) error {
	if err := r.client.SRem(ctx, key, value).Err(); err != nil {
		return c.ErrRedisOperationFailed
	}
	return nil
}

func (r *redisCacheClient) GetSetMembers(ctx context.Context, key string) ([]string, error) {
	res, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, c.ErrRedisOperationFailed
	}
	return res, nil
}

func (r *redisCacheClient) GetSetSize(ctx context.Context, key string) (int64, error) {
	res, err := r.client.SCard(ctx, key).Result()
	if err != nil {
		return 0, c.ErrRedisOperationFailed
	}
	return res, nil
}

func (r *redisCacheClient) Set(ctx context.Context, key string, ttl time.Duration, data []byte) error {
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, c.ErrCacheMiss
		}
		return nil, err
	}
	return data, nil
}

func (r *redisCacheClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCacheClient) CloseConnection() error {
	return r.client.Close()
}
