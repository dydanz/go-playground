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

// Reference : ~/server/migrations/000007_create_event_log_table.up.sql
// EventLogType represents the possible types of events
type EventLogType string

const (
	TransactionCreated   EventLogType = "transaction_created"
	ProgramIDCreated     EventLogType = "program_id_created"
	ProgramIDUpdated     EventLogType = "program_id_updated"
	PointsEarned         EventLogType = "points_earned"
	PointsRedeemed       EventLogType = "points_redeemed"
	PointsBalanceUpdated EventLogType = "points_balance_updated"
	RewardRedeemed       EventLogType = "reward_redeemed"
	UserCreated          EventLogType = "user_created"
	UserUpdated          EventLogType = "user_updated"
	MerchantCreated      EventLogType = "merchant_created"
	MerchantUpdated      EventLogType = "merchant_updated"
	ProgramCreated       EventLogType = "program_created"
	ProgramUpdated       EventLogType = "program_updated"
	ProgramRuleCreated   EventLogType = "program_rule_created"
	ProgramRuleUpdated   EventLogType = "program_rule_updated"
)

// Reference : ~/server/migrations/000007_create_event_log_table.up.sql
// EventLogActorType represents the possible types of actors
type EventLogActorType string

const (
	ClientActorType       EventLogActorType = "client"
	MerchantActorType     EventLogActorType = "merchant"
	MerchantUserActorType EventLogActorType = "merchant_user"
	SuperAdminActorType   EventLogActorType = "superadmin"
)
