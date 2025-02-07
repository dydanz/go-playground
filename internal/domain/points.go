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
	UserID    string `json:"user_id" binding:"required"`
	RewardID  string `json:"reward_id" binding:"required"`
	ProgramID string `json:"program_id" binding:"required"`
}

type UpdateRedemptionRequest struct {
	Status string `json:"status" binding:"required,oneof=completed pending failed canceled"`
}

type RedemptionService interface {
	Create(req *CreateRedemptionRequest) (*Redemption, error)
	GetByID(id string) (*Redemption, error)
	GetByUserID(userID string) ([]Redemption, error)
	Update(id string, req *UpdateRedemptionRequest) (*Redemption, error)
}

type CreateRewardRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	PointsRequired int    `json:"points_required" binding:"required,gt=0"`
	IsActive       bool   `json:"is_active"`
	Quantity       *int   `json:"quantity,omitempty"`
}

type UpdateRewardRequest struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	PointsRequired *int   `json:"points_required,omitempty"`
	IsActive       *bool  `json:"is_active,omitempty"`
	Quantity       *int   `json:"quantity,omitempty"`
}

type RewardsService interface {
	Create(req *CreateRewardRequest) (*Reward, error)
	GetByID(id string) (*Reward, error)
	GetAll(activeOnly bool) ([]Reward, error)
	Update(id string, req *UpdateRewardRequest) (*Reward, error)
	Delete(id string) error
}
