package domain

import (
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
