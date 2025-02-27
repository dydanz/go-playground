package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// SessionRepository defines the methods for session management
type SessionRepository interface {
	StoreSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetSession(ctx context.Context, userID string) (*Session, error)
	DeleteSession(ctx context.Context, userID string) error
	RefreshSession(ctx context.Context, userID, newToken string, expiration time.Duration) error
	DeleteAllSession(ctx context.Context) error
}

type Session struct {
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionRepository struct for actual implementation
type sessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) StoreSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	session := Session{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:userid:%s", userID)
	duration := time.Until(expiresAt)

	// Use context for the Redis command
	_, err = r.client.Set(ctx, key, sessionJSON, duration).Result()
	return err
}

func (r *sessionRepository) GetSession(ctx context.Context, userID string) (*Session, error) {
	key := fmt.Sprintf("session:userid:%s", userID)

	// Use context for the Redis command
	sessionJSON, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSession(ctx context.Context, userID string) error {
	key := fmt.Sprintf("session:userid:%s", userID)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete session for userID %s: %v", userID, err)
		return fmt.Errorf("failed to delete session: %w", err)
	}
	log.Printf("Successfully deleted session for userID %s", userID)
	return nil
}

func (r *sessionRepository) RefreshSession(ctx context.Context, userID, newToken string, expiration time.Duration) error {
	session, err := r.GetSession(ctx, userID)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.New("session not found")
	}

	session.TokenHash = newToken
	session.ExpiresAt = time.Now().Add(expiration)

	// Store new session and delete old one in a transaction
	pipe := r.client.TxPipeline()
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	pipe.Set(ctx, "session:userid:"+newToken, sessionJSON, expiration)
	pipe.Del(ctx, "session:userid:"+userID)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *sessionRepository) DeleteAllSession(ctx context.Context) error {
	// Get all keys with pattern "session:*"
	iter := r.client.Scan(ctx, 0, "session:*", 0).Iterator()

	// Create a pipeline for batch deletion
	pipe := r.client.Pipeline()

	// Iterate through all matching keys and queue delete commands
	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning keys: %w", err)
	}

	// Execute all delete commands
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting sessions: %w", err)
	}

	log.Printf("Successfully deleted all sessions")
	return nil
}
