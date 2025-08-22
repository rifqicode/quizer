package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures CORS settings
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// LoggerMiddleware provides structured logging
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RecoveryMiddleware handles panics and provides graceful error responses
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic with stack trace
		log.Printf("Panic recovered: %v\n%s", recovered, debug.Stack())

		// Check if headers were already written
		if c.Writer.Written() {
			return
		}

		// Return a proper error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error occurred",
			"code":    "INTERNAL_ERROR",
		})

		c.Abort()
	})
}

// ErrorHandlerMiddleware handles errors from handlers
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log the error
			log.Printf("Request error: %v", err.Error())

			// Return appropriate error response based on error type
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "Invalid request format",
					"details": err.Error(),
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
			}
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	}
}

// TimeoutMiddleware adds timeout handling to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set a timeout context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		done := make(chan struct{})

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in request handler: %v\n%s", r, debug.Stack())
				}
				close(done)
			}()
			c.Next()
		}()

		select {
		case <-done:
			// Request completed normally
		case <-ctx.Done():
			// Request timed out
			if !c.Writer.Written() {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"success": false,
					"error":   "Request timeout",
					"code":    "TIMEOUT",
				})
				c.Abort()
			}
		}
	}
}

// generateRequestID creates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
