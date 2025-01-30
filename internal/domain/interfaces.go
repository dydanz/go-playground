package domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(id string) error
	UpdateTx(tx *sql.Tx, user *User) error
	GetRandomActiveUser() (*User, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	GetUser(id string) (*User, error)
	SetUser(user *User) error
}

type AuthRepository interface {
	CreateVerification(verification *RegistrationVerification) error
	GetVerification(userID, otp string) (*RegistrationVerification, error)
	MarkVerificationUsedTx(tx *sql.Tx, verificationID string) error
	UpdateLoginAttempts(email string, increment bool) (*LoginAttempt, error)
	CreateToken(token *AuthToken) error
	InvalidateToken(userID string) error
	GetLatestVerification(userID string) (*RegistrationVerification, error)
	TxManager
}

type TxManager interface {
	BeginTx() (*sql.Tx, error)
}

type AuthService interface {
	Register(req *RegistrationRequest) (*User, error)
	Login(req *LoginRequest) (*AuthToken, error)
	Logout(userID string, tokenHash string) error
	VerifyRegistration(req *VerificationRequest) error
	GetUserByEmail(email string) (*User, error)
	GetVerificationByUserID(userID string) (*RegistrationVerification, error)
	GetRandomActiveUser() (*User, error)
}

// PointsRepository handles points balance operations
type PointsRepository interface {
	Create(ctx context.Context, ledger *PointsLedger) error
	GetByCustomerAndProgram(ctx context.Context, customerID, programID uuid.UUID) ([]*PointsLedger, error)
	GetCurrentBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error)
	GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*PointsLedger, error)
}

// TransactionRepository handles transaction operations
type TransactionRepository interface {
	Create(tx *Transaction) error
	GetByID(id string) (*Transaction, error)
	GetByUserID(userID string) ([]Transaction, error)
	Update(tx *Transaction) error
}

// RewardsRepository handles rewards operations
type RewardsRepository interface {
	Create(reward *Reward) error
	GetByID(id string) (*Reward, error)
	GetAll(activeOnly bool) ([]Reward, error)
	Update(reward *Reward) error
	Delete(id string) error
}

// RedemptionRepository handles redemption operations
type RedemptionRepository interface {
	Create(redemption *Redemption) error
	GetByID(id string) (*Redemption, error)
	GetByUserID(userID string) ([]Redemption, error)
	Update(redemption *Redemption) error
}

type PointsServiceInterface interface {
	GetLedger(ctx context.Context, customerID, programID uuid.UUID) ([]*PointsLedger, error)
	GetBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error)
	EarnPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error
	RedeemPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error
}
