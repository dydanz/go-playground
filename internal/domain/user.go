package domain

import "time"

// UserStatus represents the possible states of a user account
type UserStatus string

const (
	UserStatusPending UserStatus = "pending" // Initial state after registration
	UserStatusActive  UserStatus = "active"  // Email verified, can login
	UserStatusLocked  UserStatus = "locked"  // Account locked due to violations
	UserStatusBanned  UserStatus = "banned"  // Account banned by admin
)

// String returns the string representation of the status
func (s UserStatus) String() string {
	switch s {
	case UserStatusPending:
		return "Pending"
	case UserStatusActive:
		return "Active"
	case UserStatusLocked:
		return "Locked"
	case UserStatusBanned:
		return "Banned"
	default:
		return "Unknown"
	}
}

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
