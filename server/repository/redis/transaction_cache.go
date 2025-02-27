package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-playground/server/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type TransactionCache struct {
	client *redis.Client
}

func NewTransactionCache(client *redis.Client) *TransactionCache {
	return &TransactionCache{client: client}
}

func (c *TransactionCache) CacheTransaction(tx *domain.Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("transaction:%s", tx.TransactionID)
	return c.client.Set(context.Background(), key, data, 24*time.Hour).Err()
}

func (c *TransactionCache) GetTransaction(id string) (*domain.Transaction, error) {
	key := fmt.Sprintf("transaction:%s", id)
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var tx domain.Transaction
	if err := json.Unmarshal([]byte(data), &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (c *TransactionCache) CacheUserTransactions(userID string, transactions []domain.Transaction) error {
	data, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("user:%s:transactions", userID)
	return c.client.Set(context.Background(), key, data, 1*time.Hour).Err()
}

func (c *TransactionCache) GetUserTransactions(userID string) ([]domain.Transaction, error) {
	key := fmt.Sprintf("user:%s:transactions", userID)
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var transactions []domain.Transaction
	if err := json.Unmarshal([]byte(data), &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *TransactionCache) InvalidateUserTransactions(userID string) error {
	key := fmt.Sprintf("user:%s:transactions", userID)
	return c.client.Del(context.Background(), key).Err()
}
