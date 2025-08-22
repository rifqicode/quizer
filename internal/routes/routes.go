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

	// Initialize services once
	userService := services.NewUserService()
	tokenService := services.NewTokenService()
	authService := services.NewAuthService(userService, tokenService)
	quizService := services.NewQuizService()

	// Initialize handlers once
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	quizHandler := handlers.NewQuizHandler(quizService)

	// Health check endpoints
	router.GET("/health", healthCheck)
	router.GET("/health/db", databaseHealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		setupAuthRoutes(v1, authHandler, authService)
		setupUserRoutes(v1, userHandler, authService)
		setupQuizRoutes(v1, quizHandler, authService)
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

func setupAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authService services.AuthService) {
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

func setupUserRoutes(rg *gin.RouterGroup, userHandler *handlers.UserHandler, authService services.AuthService) {
	users := rg.Group("/users")
	users.Use(middleware.AuthMiddleware(authService)) // auth applied to all user routes
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}

func setupQuizRoutes(rg *gin.RouterGroup, quizHandler *handlers.QuizHandler, authService services.AuthService) {
	quiz := rg.Group("/quiz")
	quiz.Use(middleware.AuthMiddleware(authService)) // auth applied to all quiz routes
	{
		// Topic routes
		quiz.GET("/topics", quizHandler.GetActiveTopics)
		quiz.GET("/topics/:id", quizHandler.GetTopicByID)
		quiz.GET("/topics/:id/stats", quizHandler.GetTopicQuestionStats)

		// Quiz session routes
		quiz.POST("/sessions", quizHandler.StartQuizSession)
		quiz.GET("/sessions/:id", quizHandler.GetQuizSession)
		quiz.PUT("/sessions/:id/abandon", quizHandler.AbandonQuizSession)

		// Question routes
		quiz.GET("/sessions/:id/next-question", quizHandler.GetNextQuestion)
		quiz.POST("/sessions/:id/answer", quizHandler.SubmitAnswer)

		// Results routes
		quiz.POST("/sessions/:id/complete", quizHandler.CompleteQuizSession)
		quiz.GET("/sessions/:id/results", quizHandler.GetQuizResults)

		// User history and stats routes
		quiz.GET("/my-sessions", quizHandler.GetUserQuizHistory)
		quiz.GET("/my-stats", quizHandler.GetUserQuizStats)
	}
}
