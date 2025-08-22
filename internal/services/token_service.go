package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"quizer/datasource"
	"quizer/internal/models"
	"time"

	"gorm.io/gorm"
)

type TokenService interface {
	GenerateToken(userID uint, userAgent, ipAddress string) (*models.AccessToken, error)
	ValidateToken(token string) (*models.AccessToken, error)
	RevokeToken(token string) error
	RevokeAllUserTokens(userID uint) error
	CleanupExpiredTokens() error
}

type tokenService struct {
	db *gorm.DB
}

func NewTokenService() TokenService {
	return &tokenService{
		db: datasource.GetDB(),
	}
}

// GenerateToken creates a new access token for a user
func (s *tokenService) GenerateToken(userID uint, userAgent, ipAddress string) (*models.AccessToken, error) {
	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set expiration time (24 hours by default)
	expiresAt := time.Now().Add(24 * time.Hour)

	accessToken := &models.AccessToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		IsRevoked: false,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Load the user relation
	if err := s.db.Preload("User").First(accessToken, accessToken.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load token with user: %w", err)
	}

	return accessToken, nil
}

// ValidateToken checks if a token is valid and returns the token info
func (s *tokenService) ValidateToken(token string) (*models.AccessToken, error) {
	var accessToken models.AccessToken

	err := s.db.Preload("User").Where("token = ? AND is_revoked = ? AND expires_at > ?",
		token, false, time.Now()).First(&accessToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired token")
		}
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Double-check validity
	if !accessToken.IsValid() {
		return nil, errors.New("token is not valid")
	}

	return &accessToken, nil
}

// RevokeToken marks a specific token as revoked
func (s *tokenService) RevokeToken(token string) error {
	result := s.db.Model(&models.AccessToken{}).
		Where("token = ? AND is_revoked = ?", token, false).
		Update("is_revoked", true)

	if result.Error != nil {
		return fmt.Errorf("failed to revoke token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found or already revoked")
	}

	return nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (s *tokenService) RevokeAllUserTokens(userID uint) error {
	result := s.db.Model(&models.AccessToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true)

	if result.Error != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", result.Error)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (s *tokenService) CleanupExpiredTokens() error {
	result := s.db.Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
		Delete(&models.AccessToken{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", result.Error)
	}

	return nil
}

// GetUserActiveTokens returns all active tokens for a user
func (s *tokenService) GetUserActiveTokens(userID uint) ([]models.AccessToken, error) {
	var tokens []models.AccessToken

	err := s.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?",
		userID, false, time.Now()).Find(&tokens).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	return tokens, nil
}
