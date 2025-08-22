package models

import (
	"time"
)

// AuthRequest represents the login request payload
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email" binding:"required,email"`
	Password string `json:"password" validate:"required,min=6" binding:"required,min=6"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100" binding:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email" binding:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=128" binding:"required,min=6,max=128"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	ExpiresAt   time.Time    `json:"expires_at"`
}

// UserResponse represents user data in responses (without password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts User model to UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
