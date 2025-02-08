package main

import (
	"fmt"
	"go-playground/internal/bootstrap"
	"go-playground/internal/config"
	"go-playground/pkg/database"
	"log"
	"time"

	_ "go-playground/internal/docs" // This is required for swagger
)

// @title           Go-Playground
// @version         1.0
// @description     Go-Playground - Random Go/Gin-Boilerplate Playground
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey UserIdAuth
// @in header
// @name X-User-Id
// @description User ID for authentication

// @Security BearerAuth
// @Security UserIdAuth
func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database connections
	dbConn := bootstrap.InitializeDatabase(cfg)
	defer dbConn.RW.Close()
	defer dbConn.RR.Close()

	// Initialize Redis
	rdb := bootstrap.InitializeRedis(cfg)
	defer rdb.Close()

	// Initialize repositories
	repos := bootstrap.InitializeRepositories(dbConn.RW, dbConn, rdb, cfg)

	// Initialize services
	services := bootstrap.InitializeServices(repos)

	// Initialize handlers
	handlers := bootstrap.InitializeHandlers(services, dbConn.RW, dbConn.RR, rdb)

	// Setup router
	r := bootstrap.SetupRouter(handlers, repos.AuthRepo, repos.SessionRepo)

	// Run migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err := database.RunMigrations(dbURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Start Cleanup User Session
	repos.SessionRepo.DeleteAllSession(rdb.Context())

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := repos.AuthRepo.CleanupExpiredAttempts(); err != nil {
				log.Printf("Failed to cleanup expired attempts: %v", err)
			}
		}
	}()

	// Start server
	r.Run(":8080")
}
