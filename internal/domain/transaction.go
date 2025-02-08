package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID     uuid.UUID  `json:"transaction_id"`
	MerchantID        uuid.UUID  `json:"merchant_id"`
	CustomerID        uuid.UUID  `json:"customer_id"`
	ProgramID         uuid.UUID  `json:"program_id"`
	TransactionType   string     `json:"transaction_type"` // purchase, refund, bonus
	TransactionAmount float64    `json:"transaction_amount"`
	TransactionDate   time.Time  `json:"transaction_date"`
	BranchID          *uuid.UUID `json:"branch_id,omitempty"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
}

type CreateTransactionRequest struct {
	MerchantID        uuid.UUID  `json:"merchant_id" binding:"required"`
	CustomerID        uuid.UUID  `json:"customer_id" binding:"required"`
	ProgramID         uuid.UUID  `json:"program_id" binding:"required"`
	TransactionType   string     `json:"transaction_type" binding:"required,oneof=purchase refund bonus"`
	TransactionAmount float64    `json:"transaction_amount" binding:"required,gt=0"`
	BranchID          *uuid.UUID `json:"branch_id,omitempty"`
	Status            string     `json:"status" binding:"required,oneof=pending completed failed cancelled"`
}

type UpdateTransactionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed failed cancelled"`
}
