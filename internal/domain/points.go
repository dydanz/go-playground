package domain

import (
	"time"

	"github.com/google/uuid"
)

type PointsLedger struct {
	LedgerID            uuid.UUID  `json:"ledger_id"`
	MerchantCustomersID uuid.UUID  `json:"merchant_customers_id"`
	ProgramID           uuid.UUID  `json:"program_id"`
	PointsEarned        int        `json:"points_earned"`
	PointsRedeemed      int        `json:"points_redeemed"`
	PointsBalance       int        `json:"points_balance"`
	TransactionID       *uuid.UUID `json:"transaction_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

type Reward struct {
	ID                uuid.UUID `json:"id"`
	ProgramID         uuid.UUID `json:"program_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	PointsRequired    int       `json:"points_required"`
	AvailableQuantity *int      `json:"available_quantity,omitempty"`
	Quantity          int       `json:"quantity"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RedemptionStatus string

const (
	RedemptionStatusCompleted RedemptionStatus = "completed"
	RedemptionStatusPending   RedemptionStatus = "pending"
	RedemptionStatusFailed    RedemptionStatus = "failed"
)

type Redemption struct {
	ID                  uuid.UUID        `json:"id"`
	MerchantCustomersID uuid.UUID        `json:"merchant_customers_id"`
	RewardID            uuid.UUID        `json:"reward_id"`
	PointsUsed          int              `json:"points_used"`
	RedemptionDate      time.Time        `json:"redemption_date"`
	Status              RedemptionStatus `json:"status"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

type PointsBalance struct {
	CustomerID string `json:"customer_id"`
	ProgramID  string `json:"program_id"`
	Balance    int    `json:"balance"`
}

type PointsTransaction struct {
	TransactionID string    `json:"transaction_id"`
	CustomerID    string    `json:"customer_id"`
	ProgramID     string    `json:"program_id"`
	Points        int       `json:"points"`
	Type          string    `json:"type"` // "earn" or "redeem"
	CreatedAt     time.Time `json:"created_at"`
}

type EarnPointsRequest struct {
	CustomerID string `json:"customer_id"`
	ProgramID  string `json:"program_id"`
	Points     int    `json:"points" binding:"required,gt=0"`
}

type RedeemPointsRequest struct {
	CustomerID string `json:"customer_id"`
	ProgramID  string `json:"program_id"`
	Points     int    `json:"points" binding:"required,gt=0"`
}

type CreateRedemptionRequest struct {
	MerchantCustomersID uuid.UUID `json:"merchant_customers_id" binding:"required"`
	RewardID            uuid.UUID `json:"reward_id" binding:"required"`
	PointsUsed          int       `json:"points_used" binding:"required,gt=0"`
}

type UpdateRedemptionRequest struct {
	Status RedemptionStatus `json:"status" binding:"required,oneof=completed pending failed"`
}


type CreateRewardRequest struct {
	ProgramID         uuid.UUID `json:"program_id" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	Description       string    `json:"description" binding:"required"`
	PointsRequired    int       `json:"points_required" binding:"required,gt=0"`
	AvailableQuantity *int      `json:"available_quantity,omitempty"`
	Quantity          int       `json:"quantity"`
	IsActive          bool      `json:"is_active"`
}

type UpdateRewardRequest struct {
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	PointsRequired    *int   `json:"points_required,omitempty"`
	AvailableQuantity *int   `json:"available_quantity,omitempty"`
	Quantity          *int   `json:"quantity,omitempty"`
	IsActive          *bool  `json:"is_active,omitempty"`
}
