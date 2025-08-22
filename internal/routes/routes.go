package routes

import (
	"net/http"
	"quizer/datasource"
	"quizer/internal/handlers"
	"quizer/internal/middleware"
	"quizer/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.New()

	// Core middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.TimeoutMiddleware(30 * time.Second)) // 30 second timeout

	// Health check endpoints
	router.GET("/health", healthCheck)
	router.GET("/health/db", databaseHealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		setupAuthRoutes(v1)
		setupUserRoutes(v1)
	}

	return router
}

// healthCheck provides basic health status
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"message":   "Server is running",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// databaseHealthCheck checks database connectivity
func databaseHealthCheck(c *gin.Context) {
	if datasource.CheckDB() {
		c.JSON(http.StatusOK, gin.H{
			"status":    "OK",
			"message":   "Database is healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "ERROR",
			"message":   "Database is unhealthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func setupAuthRoutes(rg *gin.RouterGroup) {
	userService := services.NewUserService()
	tokenService := services.NewTokenService()
	authService := services.NewAuthService(userService, tokenService)
	authHandler := handlers.NewAuthHandler(authService, userService)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)

		// Protected routes
		auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
	}
}

func setupUserRoutes(rg *gin.RouterGroup) {
	userService := services.NewUserService()
	tokenService := services.NewTokenService()
	authService := services.NewAuthService(userService, tokenService)
	userHandler := handlers.NewUserHandler(userService)

	users := rg.Group("/users")
	// Apply authentication middleware to all user routes
	users.Use(middleware.AuthMiddleware(authService))
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
