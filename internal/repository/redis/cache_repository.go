package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-cursor/internal/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *CacheRepository) Get(key string, dest interface{}) error {
	ctx := context.Background()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *CacheRepository) GetUser(id string) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", id)
	data, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *CacheRepository) SetUser(user *domain.User) error {
	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), key, data, 24*time.Hour).Err()
}
