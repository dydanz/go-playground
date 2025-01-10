package main

import (
	"go-cursor/internal/config"
	"go-cursor/internal/handler"
	"go-cursor/internal/repository/postgres"
	"go-cursor/internal/repository/redis"
	"go-cursor/internal/service"
	"go-cursor/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize PostgreSQL
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	rdb := database.NewRedisConnection(cfg)
	defer rdb.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	cacheRepo := redis.NewCacheRepository(rdb)

	// Initialize services
	userService := service.NewUserService(userRepo, cacheRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Initialize Gin router
	r := gin.Default()

	// Routes
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.POST("/", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	// Start server
	r.Run(":8080")
}
