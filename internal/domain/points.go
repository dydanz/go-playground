package domain

import (
	"time"

	"github.com/google/uuid"
)

type PointsLedger struct {
	LedgerID       uuid.UUID  `json:"ledger_id"`
	CustomerID     uuid.UUID  `json:"customer_id"`
	ProgramID      uuid.UUID  `json:"program_id"`
	PointsEarned   int        `json:"points_earned"`
	PointsRedeemed int        `json:"points_redeemed"`
	PointsBalance  int        `json:"points_balance"`
	TransactionID  *uuid.UUID `json:"transaction_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type Reward struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	PointsRequired int       `json:"points_required"`
	IsActive       bool      `json:"is_active"`
	Quantity       *int      `json:"quantity,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Redemption struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RewardID   string    `json:"reward_id"`
	ProgramID  string    `json:"program_id"`
	Status     string    `json:"status"` // "completed", "pending", "failed", "canceled"
	RedeemedAt time.Time `json:"redeemed_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
