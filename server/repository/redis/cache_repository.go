package redis

import (
	"context"
	"encoding/json"
	"go-playground/server/domain"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *CacheRepository) Get(ctx context.Context, key string) (interface{}, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *CacheRepository) SetUser(user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.Set(context.Background(), "user:"+user.ID, data)
}

func (r *CacheRepository) GetUser(id string) (*domain.User, error) {
	data, err := r.Get(context.Background(), "user:"+id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data.(string)), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
