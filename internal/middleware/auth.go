package middleware

import (
	"net/http"
	"quizer/internal/services"
	"quizer/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates token authentication middleware
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		user, err := authService.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Account is deactivated")
			c.Abort()
			return
		}

		// Store user information in context
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user", user)

		c.Next()
	}
}
