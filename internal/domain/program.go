package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProgramRule struct {
	ID             uuid.UUID  `json:"id"`
	ProgramID      uuid.UUID  `json:"program_id"`
	RuleName       string     `json:"rule_name"`
	ConditionType  string     `json:"condition_type"`
	ConditionValue string     `json:"condition_value"`
	Multiplier     float64    `json:"multiplier"`
	PointsAwarded  int        `json:"points_awarded"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateProgramRuleRequest struct {
	ProgramID      uuid.UUID  `json:"program_id" binding:"required"`
	RuleName       string     `json:"rule_name" binding:"required"`
	ConditionType  string     `json:"condition_type" binding:"required"`
	ConditionValue string     `json:"condition_value" binding:"required"`
	Multiplier     float64    `json:"multiplier" binding:"required,gt=0"`
	PointsAwarded  int        `json:"points_awarded" binding:"required,gte=0"`
	EffectiveFrom  time.Time  `json:"effective_from" binding:"required"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
}

type UpdateProgramRuleRequest struct {
	RuleName       string     `json:"rule_name,omitempty"`
	ConditionType  string     `json:"condition_type,omitempty"`
	ConditionValue string     `json:"condition_value,omitempty"`
	Multiplier     *float64   `json:"multiplier,omitempty"`
	PointsAwarded  *int       `json:"points_awarded,omitempty"`
	EffectiveFrom  *time.Time `json:"effective_from,omitempty"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
}

type CreateProgramResponse struct {
	ProgramID         uuid.UUID `json:"program_id"`
	MerchantID        uuid.UUID `json:"merchant_id"`
	ProgramName       string    `json:"program_name"`
	PointCurrencyName string    `json:"point_currency_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Program struct {
	ID                uuid.UUID `json:"program_id"`
	MerchantID        uuid.UUID `json:"merchant_id"`
	UserID            uuid.UUID `json:"user_id"`
	ProgramName       string    `json:"program_name"`
	PointCurrencyName string    `json:"point_currency_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateProgramRequest struct {
	MerchantID        uuid.UUID `json:"merchant_id" binding:"required"`
	UserID            uuid.UUID `json:"user_id" binding:"required"`
	ProgramName       string    `json:"program_name" binding:"required"`
	PointCurrencyName string    `json:"point_currency_name" binding:"required"`
}

type UpdateProgramRequest struct {
	ProgramName       string `json:"program_name,omitempty"`
	PointCurrencyName string `json:"point_currency_name,omitempty"`
}
