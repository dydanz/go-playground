package domain

import "time"

type User struct {
	ID        string    `json:"id"` // UUID
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hidden from JSON responses
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
