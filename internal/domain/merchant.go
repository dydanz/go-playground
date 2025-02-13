package domain

import (
	"context"
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

type MerchantRepository interface {
	Create(ctx context.Context, merchant *Merchant) (*Merchant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Merchant, error)
	GetAll(ctx context.Context) ([]*Merchant, error)
	Update(ctx context.Context, merchant *Merchant) error
	Delete(ctx context.Context, id uuid.UUID) error
}
