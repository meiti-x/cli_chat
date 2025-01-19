package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/meiti-x/snapp_chal/config"
)

func MustInitRedis(conf *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
}
