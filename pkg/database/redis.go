package database

import (
	"fmt"
	"go-cursor/internal/config"

	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
}
