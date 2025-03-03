package redis

import (
	"context"
	"encoding/json"
	"go-playground/pkg/logging"
	"go-playground/server/domain"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type CacheRepository struct {
	client *redis.Client
	logger zerolog.Logger
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client,
		logger: logging.GetLogger(),
	}
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
		r.logger.Error().
			Err(err).
			Msg("Failed to marshal user")
		return err
	}
	return r.Set(context.Background(), "user:"+user.ID, data)
}

func (r *CacheRepository) GetUser(id string) (*domain.User, error) {
	data, err := r.Get(context.Background(), "user:"+id)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to get user")
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data.(string)), &user); err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to unmarshal user")
		return nil, err
	}
	return &user, nil
}
