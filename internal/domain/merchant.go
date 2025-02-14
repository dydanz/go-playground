package domain

import (
	"time"

	"github.com/google/uuid"
)

type MerchantType string

const (
	MerchantTypeBank       MerchantType = "bank"
	MerchantTypeEcommerce  MerchantType = "e-commerce"
	MerchantTypeRepairShop MerchantType = "repair_shop"
)

type Merchant struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	Name      string       `json:"merchant_name"`
	Type      MerchantType `json:"merchant_type"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateMerchantRequest struct {
	UserID uuid.UUID    `json:"user_id" binding:"required"`
	Name   string       `json:"merchant_name" binding:"required"`
	Type   MerchantType `json:"merchant_type" binding:"required,oneof=bank e-commerce repair_shop"`
}

type UpdateMerchantRequest struct {
	Name string       `json:"merchant_name" binding:"required"`
	Type MerchantType `json:"merchant_type" binding:"required,oneof=bank e-commerce repair_shop"`
}

// Pagination represents pagination metadata
type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
}

// PaginatedMerchants represents a paginated list of merchants
type PaginatedMerchants struct {
	Merchants  []*Merchant `json:"merchants"`
	Pagination Pagination  `json:"pagination"`
}
