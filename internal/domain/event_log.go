package domain

import (
	"context"
	"time"
)

// EventLog represents a log entry for events
type EventLog struct {
	ID             string                 `json:"id"`
	EventType      string                 `json:"event_type"`
	UserID         string                 `json:"user_id"`
	Details        map[string]interface{} `json:"details"`
	EventTimestamp time.Time              `json:"event_timestamp"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ReferenceID    *string                `json:"reference_id,omitempty"`
}

type EventLogRepository interface {
	Create(ctx context.Context, eventLog *EventLog) error
	GetByID(id string) (*EventLog, error)
	GetByUserID(userID string) ([]EventLog, error)
	Update(eventLog *EventLog) error
	Delete(id string) error
	GetByReferenceID(referenceID string) (*EventLog, error)
}
