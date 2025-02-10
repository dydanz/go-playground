package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// MerchantCustomer represents a customer of a merchant
type MerchantCustomer struct {
	ID         uuid.UUID `json:"id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MerchantCustomersRepository defines the interface for merchant customer data operations
type MerchantCustomersRepository interface {
	Create(ctx context.Context, customer *MerchantCustomer) error
	GetByID(ctx context.Context, id uuid.UUID) (*MerchantCustomer, error)
	GetByEmail(ctx context.Context, email string) (*MerchantCustomer, error)
	GetByPhone(ctx context.Context, phone string) (*MerchantCustomer, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MerchantCustomer, error)
	Update(ctx context.Context, customer *MerchantCustomer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// MerchantCustomersService defines the interface for merchant customer business logic
type MerchantCustomersService interface {
	Create(ctx context.Context, req *CreateMerchantCustomerRequest) (*MerchantCustomer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*MerchantCustomer, error)
	GetByEmail(ctx context.Context, email string) (*MerchantCustomer, error)
	GetByPhone(ctx context.Context, phone string) (*MerchantCustomer, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MerchantCustomer, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateMerchantCustomerRequest) (*MerchantCustomer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ValidateCredentials(ctx context.Context, email, password string) (*MerchantCustomer, error)
}

// CreateMerchantCustomerRequest represents the request to create a new merchant customer
type CreateMerchantCustomerRequest struct {
	MerchantID uuid.UUID `json:"merchant_id" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required,min=6"`
	Name       string    `json:"name" validate:"required"`
	Phone      string    `json:"phone" validate:"required"`
}

// UpdateMerchantCustomerRequest represents the request to update an existing merchant customer
type UpdateMerchantCustomerRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

// CustomerLoginRequest represents the login request for merchant customers
type CustomerLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
