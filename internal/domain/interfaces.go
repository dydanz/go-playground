package domain

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(id string) error
	UpdateTx(tx *sql.Tx, user *User) error
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
	TxManager
}

type TxManager interface {
	BeginTx() (*sql.Tx, error)
}

type AuthService interface {
	Register(req *RegistrationRequest) (*User, error)
	Login(req *LoginRequest) (*AuthToken, error)
	Logout(userID string) error
	VerifyRegistration(req *VerificationRequest) error
}
