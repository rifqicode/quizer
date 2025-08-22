package handlers

import (
	"log"
	"net/http"
	"quizer/internal/models"
	"quizer/internal/services"
	"quizer/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param register body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Register: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	var req models.RegisterRequest

	if err := utils.BindJSON(c, &req); err != nil {
		return // BindJSON already sends error response
	}

	// Additional validation
	if req.Name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Name is required")
		return
	}

	// Get user agent and IP address
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	authResponse, err := h.authService.Register(&req, userAgent, ipAddress)
	if err != nil {
		log.Printf("Error registering user: %v", err)

		// Check for specific error types
		if err.Error() == "user with this email already exists" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", authResponse)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.AuthRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Login: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	var req models.AuthRequest

	if err := utils.BindJSON(c, &req); err != nil {
		return // BindJSON already sends error response
	}

	// Get user agent and IP address
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	authResponse, err := h.authService.Login(&req, userAgent, ipAddress)
	if err != nil {
		log.Printf("Error logging in user: %v", err)

		// Don't expose specific error details for security
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", authResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token from existing valid token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in RefreshToken: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
		return
	}

	// Remove "Bearer " prefix
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}
	token := authHeader[7:] // Skip "Bearer "

	// Get user agent and IP address
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	authResponse, err := h.authService.RefreshToken(token, userAgent, ipAddress)
	if err != nil {
		log.Printf("Error refreshing token: %v", err)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", authResponse)
}

// Me godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Me: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", user.ToUserResponse())
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate the access token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Logout: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
		return
	}

	// Remove "Bearer " prefix
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}
	token := authHeader[7:] // Skip "Bearer "

	// Revoke the token
	err := h.authService.Logout(token)
	if err != nil {
		log.Printf("Error during logout: %v", err)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logout successful", map[string]string{
		"message": "Token has been successfully revoked",
	})
}
