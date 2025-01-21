package domain

import "time"

type PointsBalance struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TotalPoints int       `json:"total_points"`
	LastUpdated time.Time `json:"last_updated"`
}

type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	TransactionType string    `json:"transaction_type"` // "earn" or "redeem"
	Points          int       `json:"points"`
	Description     string    `json:"description"`
	Status          string    `json:"status"` // "completed", "pending", "canceled"
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
	Status     string    `json:"status"` // "completed", "pending", "failed", "canceled"
	RedeemedAt time.Time `json:"redeemed_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
