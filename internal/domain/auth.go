package domain

import (
	"time"
)

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type VerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthToken struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	TokenHash  string    `json:"-"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
}

type LoginAttempt struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	AttemptCount int       `json:"attempt_count"`
	LastAttempt  time.Time `json:"last_attempt"`
	LockedUntil  time.Time `json:"locked_until,omitempty"`
}

type RegistrationVerification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UsedAt    time.Time `json:"used_at,omitempty"`
}

// LoginResponse represents the response from a successful login
type LoginResponse struct {
	Token     string    `json:"token" example:"Bearer eyJhbGciOiJ..."` // Token with Bearer prefix
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
}
