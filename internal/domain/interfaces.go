package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	UpdateTx(ctx context.Context, tx *sql.Tx, user *User) error
	GetRandomActiveUser(ctx context.Context) (*User, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	GetUser(id string) (*User, error)
	SetUser(user *User) error
}

type AuthRepository interface {
	CreateVerification(ctx context.Context, verification *RegistrationVerification) error
	GetVerification(ctx context.Context, userID, otp string) (*RegistrationVerification, error)
	MarkVerificationUsedTx(ctx context.Context, tx *sql.Tx, verificationID string) error
	UpdateLoginAttempts(ctx context.Context, email string, increment bool) (*LoginAttempt, error)
	CreateToken(ctx context.Context, token *AuthToken) error
	InvalidateToken(ctx context.Context, userID string) error
	GetLatestVerification(ctx context.Context, userID string) (*RegistrationVerification, error)
	TxManager
}

type TxManager interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type AuthService interface {
	Register(ctx context.Context, req *RegistrationRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthToken, error)
	Logout(ctx context.Context, userID string, tokenHash string) error
	VerifyRegistration(ctx context.Context, req *VerificationRequest) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetVerificationByUserID(ctx context.Context, userID string) (*RegistrationVerification, error)
	GetRandomActiveUser(ctx context.Context) (*User, error)
}

// PointsRepository handles points balance operations
type PointsRepository interface {
	Create(ctx context.Context, ledger *PointsLedger) (*PointsLedger, error)
	GetByCustomerAndProgram(ctx context.Context, customerID, programID uuid.UUID) ([]*PointsLedger, error)
	GetCurrentBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error)
	GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*PointsLedger, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) (*Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Transaction, error)
	GetByCustomerIDWithPagination(ctx context.Context, customerID uuid.UUID, offset, limit int) ([]*Transaction, int64, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// RewardsRepository handles rewards operations
type RewardsRepository interface {
	Create(ctx context.Context, reward *Reward) (*Reward, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Reward, error)
	Update(ctx context.Context, reward *Reward) (*Reward, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, activeOnly bool) ([]Reward, error)
}

// RedemptionRepository handles redemption operations
type RedemptionRepository interface {
	Create(ctx context.Context, redemption *Redemption) ([]*Redemption, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Redemption, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Redemption, error)
	Update(ctx context.Context, redemption *Redemption) error
}

type ProgramRepository interface {
	Create(ctx context.Context, program *Program) (*Program, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Program, error)
	GetAll(ctx context.Context) ([]*Program, error)
	Update(ctx context.Context, program *Program) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Program, error)
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
	GetLedger(ctx context.Context, customerID uuid.UUID, programID uuid.UUID) ([]*PointsLedger, error)
	GetBalance(ctx context.Context, customerID uuid.UUID, programID uuid.UUID) (*PointsBalance, error)
	EarnPoints(ctx context.Context, req *PointsTransaction) (*PointsTransaction, error)
	RedeemPoints(ctx context.Context, req *PointsTransaction) (*PointsTransaction, error)
}

type ProgramService interface {
	Create(ctx context.Context, req *CreateProgramRequest) (*Program, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Program, error)
	GetAll(ctx context.Context) ([]*Program, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateProgramRequest) (*Program, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Program, error)
}

type MerchantService interface {
	Create(ctx context.Context, req *CreateMerchantRequest) (*Merchant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	GetAll(ctx context.Context) ([]*Merchant, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateMerchantRequest) (*Merchant, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProgramRulesService interface {
	Create(req *CreateProgramRuleRequest) (*ProgramRule, error)
	GetByID(id string) (*ProgramRule, error)
	GetByProgramID(programID string) ([]*ProgramRule, error)
	Update(id string, req *UpdateProgramRuleRequest) (*ProgramRule, error)
	Delete(id string) error
	GetActiveRules(programID string) ([]*ProgramRule, error)
}

type RewardsService interface {
	Create(ctx context.Context, req *CreateRewardRequest) (*Reward, error)
	GetByID(ctx context.Context, id string) (*Reward, error)
	Update(ctx context.Context, id string, req *UpdateRewardRequest) (*Reward, error)
	Delete(ctx context.Context, id string) error
}

type RedemptionService interface {
	Create(req *CreateRedemptionRequest) (*Redemption, error)
	GetByID(id string) (*Redemption, error)
	GetByUserID(userID string) ([]Redemption, error)
	Update(id string, req *UpdateRedemptionRequest) (*Redemption, error)
}

type EventLogRepository interface {
	Create(ctx context.Context, eventLog *EventLog) error
	GetByID(id string) (*EventLog, error)
	GetByUserID(userID string) ([]EventLog, error)
	GetByReferenceID(referenceID string) (*EventLog, error)
}

// TransactionRepository handles transaction operations
type TransactionService interface {
	Create(ctx context.Context, req *CreateTransactionRequest) (*Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Transaction, error)
	GetByCustomerIDWithPagination(ctx context.Context, customerID uuid.UUID, offset, limit int) ([]*Transaction, int64, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Transaction, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
