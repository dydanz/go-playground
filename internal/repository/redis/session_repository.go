package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	client *redis.Client
}

type Session struct {
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) StoreSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
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

func (r *SessionRepository) GetSession(ctx context.Context, userID string) (*Session, error) {
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

func (r *SessionRepository) DeleteSession(ctx context.Context, userid string) error {
	return r.client.Del(ctx, "session:userid"+userid).Err()
}

func (r *SessionRepository) RefreshSession(ctx context.Context, userid, newToken string, expiration time.Duration) error {
	session, err := r.GetSession(ctx, userid)
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
	pipe.Del(ctx, "session:userid:"+userid)

	_, err = pipe.Exec(ctx)
	return err
}
