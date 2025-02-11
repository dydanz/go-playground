package domain

import (
	"time"
)

// EventLog represents a log entry for events
type EventLog struct {
	ID             string                 `json:"id"`
	EventType      string                 `json:"event_type"`
	ActorID        string                 `json:"actor_id"`
	ActorType      string                 `json:"actor_type"`
	Details        map[string]interface{} `json:"details"`
	EventTimestamp time.Time              `json:"event_timestamp"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ReferenceID    *string                `json:"reference_id,omitempty"`
}

// Reference : ~/internal/migrations/000007_create_event_log_table.up.sql
// EventLogType represents the possible types of events
type EventLogType string

const (
	TransactionCreated   EventLogType = "transaction_created"    // Initial state after registration
	ProgramIDCreated     EventLogType = "program_id_created"     // Email verified, can login
	ProgramRuleCreated   EventLogType = "program_rule_created"   // Account locked due to violations
	PointsEarned         EventLogType = "points_earned"          // Account banned by admin
	PointsRedeemed       EventLogType = "points_redeemed"        // Account banned by admin
	PointsBalanceUpdated EventLogType = "points_balance_updated" // Account banned by admin
	RewardRedeemed       EventLogType = "reward_redeemed"        // Account banned by admin
)

// Reference : ~/internal/migrations/000007_create_event_log_table.up.sql
// EventLogActorType represents the possible types of actors
type EventLogActorType string

const (
	ClientActorType       EventLogActorType = "client"
	MerchantActorType     EventLogActorType = "merchant"
	MerchantUserActorType EventLogActorType = "merchant_user"
	SuperAdminActorType   EventLogActorType = "superadmin"
)
