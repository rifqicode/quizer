package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quizer/datasource"
	"quizer/internal/config"
	"quizer/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration first
	config.LoadConfig()

	// Connect to database with error handling
	db := datasource.ConnectDatabase()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// Ensure database connection is closed on exit
	defer func() {
		if err := datasource.CloseDB(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.Mode)

	// Setup routes
	router := routes.SetupRoutes()

	// Configure HTTP server
	srv := &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		Handler:      router,
		ReadTimeout:  config.AppConfig.Server.ReadTimeout,
		WriteTimeout: config.AppConfig.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
		log.Printf("Server mode: %s", config.AppConfig.Server.Mode)
		log.Printf("Health check: http://localhost:%s/health", config.AppConfig.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Setup graceful shutdown
	setupGracefulShutdown(ctx, srv)
}

// setupGracefulShutdown handles graceful server shutdown
func setupGracefulShutdown(ctx context.Context, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Server is shutting down...")

	// Create a context with timeout for shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return
	}

	log.Println("Server exited gracefully")
}
