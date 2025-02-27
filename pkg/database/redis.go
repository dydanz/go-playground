package database

import (
	"fmt"
	"go-playground/server/config"

	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(cfg *config.Config) *redis.Client {
	options := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		DB:   0,
	}

	// Only set password if it's not empty
	if cfg.RedisPassword != "" {
		options.Password = cfg.RedisPassword
	}

	return redis.NewClient(options)
}
