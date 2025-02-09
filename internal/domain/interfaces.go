package domain

import (
	"context"
	"database/sql"
	"time"

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
type TransactionService interface {
	Create(req *CreateTransactionRequest) (*Transaction, error)
	GetByID(id string) (*Transaction, error)
	GetByCustomerID(customerID string) ([]*Transaction, error)
	GetByCustomerIDWithPagination(customerID string, offset, limit int) ([]*Transaction, int64, error)
	GetByMerchantID(merchantID string) ([]*Transaction, error)
	UpdateStatus(id string, status string) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Transaction, error)
	GetByCustomerIDWithPagination(ctx context.Context, customerID uuid.UUID, offset, limit int) ([]*Transaction, int64, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// RewardsRepository handles rewards operations
type RewardsRepository interface {
	Create(reward *Reward) error
	GetByID(id uuid.UUID) (*Reward, error)
	Update(reward *Reward) error
	Delete(id uuid.UUID) error
}

// RedemptionRepository handles redemption operations
type RedemptionRepository interface {
	Create(ctx context.Context, redemption *Redemption) error
	GetByID(ctx context.Context, id uuid.UUID) (*Redemption, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Redemption, error)
	Update(ctx context.Context, redemption *Redemption) error
}

type PointsServiceInterface interface {
	GetLedger(ctx context.Context, customerID, programID uuid.UUID) ([]*PointsLedger, error)
	GetBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error)
	EarnPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error
	RedeemPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error
}

type ProgramRuleRepository interface {
	Create(ctx context.Context, rule *ProgramRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*ProgramRule, error)
	GetByProgramID(ctx context.Context, programID uuid.UUID) ([]*ProgramRule, error)
	Update(ctx context.Context, rule *ProgramRule) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveRules(ctx context.Context, programID uuid.UUID, timestamp time.Time) ([]*ProgramRule, error)
}

type UserService interface {
	Create(req *CreateUserRequest) (*User, error)
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	Update(id string, req *UpdateUserRequest) (*User, error)
	Delete(id string) error
}

type PointsService interface {
	GetLedger(customerID string, programID string) (*PointsLedger, error)
	GetBalance(customerID string, programID string) (*PointsBalance, error)
	EarnPoints(req *EarnPointsRequest) (*PointsTransaction, error)
	RedeemPoints(req *RedeemPointsRequest) (*PointsTransaction, error)
}

type ProgramService interface {
	Create(req *CreateProgramRequest) (*Program, error)
	GetByID(id string) (*Program, error)
	GetAll() ([]*Program, error)
	Update(id string, req *UpdateProgramRequest) (*Program, error)
	Delete(id string) error
	GetByMerchantID(merchantID string) ([]*Program, error)
}

type MerchantService interface {
	Create(req *CreateMerchantRequest) (*Merchant, error)
	GetByID(id uuid.UUID) (*Merchant, error)
	GetAll() ([]*Merchant, error)
	Update(id uuid.UUID, req *UpdateMerchantRequest) (*Merchant, error)
	Delete(id uuid.UUID) error
}

type ProgramRulesService interface {
	Create(req *CreateProgramRuleRequest) (*ProgramRule, error)
	GetByID(id string) (*ProgramRule, error)
	GetByProgramID(programID string) ([]*ProgramRule, error)
	Update(id string, req *UpdateProgramRuleRequest) (*ProgramRule, error)
	Delete(id string) error
	GetActiveRules(programID string) ([]*ProgramRule, error)
}

type ProgramRepository interface {
	Create(program *Program) error
	GetByID(id string) (*Program, error)
	GetAll() ([]*Program, error)
	Update(program *Program) error
	Delete(id string) error
	GetByMerchantID(merchantID string) ([]*Program, error)
}

type RewardsService interface {
	Create(req *CreateRewardRequest) (*Reward, error)
	GetByID(id string) (*Reward, error)
	Update(id string, req *UpdateRewardRequest) (*Reward, error)
	Delete(id string) error
}

type RedemptionService interface {
	Create(req *CreateRedemptionRequest) (*Redemption, error)
	GetByID(id string) (*Redemption, error)
	GetByUserID(userID string) ([]Redemption, error)
	Update(id string, req *UpdateRedemptionRequest) (*Redemption, error)
}
