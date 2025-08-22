package models

import (
	"time"
)

type AccessToken struct {
	BaseModel
	Token     string    `json:"token" gorm:"unique;not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false;index"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	IPAddress string    `json:"ip_address" gorm:"type:varchar(45)"`
}

// IsValid checks if the token is still valid
func (at *AccessToken) IsValid() bool {
	return !at.IsRevoked && time.Now().Before(at.ExpiresAt)
}

// Revoke marks the token as revoked
func (at *AccessToken) Revoke() {
	at.IsRevoked = true
}
