package services

import (
	"errors"
	"fmt"
	"quizer/datasource"
	"quizer/internal/models"
	"quizer/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *models.RegisterRequest, userAgent, ipAddress string) (*models.AuthResponse, error)
	Login(req *models.AuthRequest, userAgent, ipAddress string) (*models.AuthResponse, error)
	Logout(token string) error
	ValidateToken(tokenString string) (*models.User, error)
	RefreshToken(tokenString string, userAgent, ipAddress string) (*models.AuthResponse, error)
}

type authService struct {
	db           *gorm.DB
	userService  UserService
	tokenService TokenService
}

func NewAuthService(userService UserService, tokenService TokenService) AuthService {
	return &authService{
		db:           datasource.GetDB(),
		userService:  userService,
		tokenService: tokenService,
	}
}

// Register creates a new user account
func (s *authService) Register(req *models.RegisterRequest, userAgent, ipAddress string) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userService.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Validate password strength
	if passwordErrors := utils.ValidatePasswordStrength(req.Password); len(passwordErrors) > 0 {
		return nil, fmt.Errorf("password validation failed: %v", passwordErrors)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		IsActive: true,
	}

	if err := s.userService.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateToken(user.ID, userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		User:        user.ToUserResponse(),
		AccessToken: accessToken.Token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(accessToken.ExpiresAt).Seconds()),
		ExpiresAt:   accessToken.ExpiresAt,
	}, nil
}

// Login authenticates a user and returns an access token
func (s *authService) Login(req *models.AuthRequest, userAgent, ipAddress string) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userService.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateToken(user.ID, userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		User:        user.ToUserResponse(),
		AccessToken: accessToken.Token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(accessToken.ExpiresAt).Seconds()),
		ExpiresAt:   accessToken.ExpiresAt,
	}, nil
}

// Logout invalidates the provided token
func (s *authService) Logout(token string) error {
	return s.tokenService.RevokeToken(token)
}

// ValidateToken validates an access token and returns the user
func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	accessToken, err := s.tokenService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Return the user from the access token
	return &accessToken.User, nil
}

// RefreshToken generates a new token from an existing valid token
func (s *authService) RefreshToken(tokenString string, userAgent, ipAddress string) (*models.AuthResponse, error) {
	// Validate the current token
	accessToken, err := s.tokenService.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	// Revoke the old token
	if err := s.tokenService.RevokeToken(tokenString); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate a new token
	newAccessToken, err := s.tokenService.GenerateToken(accessToken.UserID, userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token: %w", err)
	}

	return &models.AuthResponse{
		User:        accessToken.User.ToUserResponse(),
		AccessToken: newAccessToken.Token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(newAccessToken.ExpiresAt).Seconds()),
		ExpiresAt:   newAccessToken.ExpiresAt,
	}, nil
}
