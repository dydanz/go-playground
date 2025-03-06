package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type TransactionCache struct {
	client *redis.Client
	logger zerolog.Logger
}

func NewTransactionCache(client *redis.Client) *TransactionCache {
	return &TransactionCache{client: client,
		logger: logging.GetLogger(),
	}
}

func (c *TransactionCache) CacheTransaction(tx *domain.Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to marshal transaction")
		return err
	}

	key := fmt.Sprintf("transaction:%s", tx.TransactionID)
	return c.client.Set(context.Background(), key, data, 24*time.Hour).Err()
}

func (c *TransactionCache) GetTransaction(id string) (*domain.Transaction, error) {
	key := fmt.Sprintf("transaction:%s", id)
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		c.logger.Error().
			Msg("Transaction not found")
		return nil, nil
	}
	if err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to get transaction")
		return nil, err
	}

	var tx domain.Transaction
	if err := json.Unmarshal([]byte(data), &tx); err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to unmarshal transaction")
		return nil, err
	}
	return &tx, nil
}

func (c *TransactionCache) CacheUserTransactions(userID string, transactions []domain.Transaction) error {
	data, err := json.Marshal(transactions)
	if err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to marshal transactions")
		return err
	}

	key := fmt.Sprintf("user:%s:transactions", userID)
	return c.client.Set(context.Background(), key, data, 1*time.Hour).Err()
}

func (c *TransactionCache) GetUserTransactions(userID string) ([]domain.Transaction, error) {
	key := fmt.Sprintf("user:%s:transactions", userID)
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		c.logger.Error().
			Msg("Transactions not found")
		return nil, nil
	}
	if err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to get transactions")
		return nil, err
	}

	var transactions []domain.Transaction
	if err := json.Unmarshal([]byte(data), &transactions); err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to unmarshal transactions")
		return nil, err
	}
	return transactions, nil
}

func (c *TransactionCache) InvalidateUserTransactions(userID string) error {
	key := fmt.Sprintf("user:%s:transactions", userID)
	err := c.client.Del(context.Background(), key).Err()
	if err != nil {
		c.logger.Error().
			Err(err).
			Msg("Failed to invalidate user transactions")
		return err
	}
	return nil
}
