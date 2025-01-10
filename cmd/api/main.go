package main

import (
	"go-cursor/internal/config"
	"go-cursor/internal/handler"
	"go-cursor/internal/repository/postgres"
	"go-cursor/internal/repository/redis"
	"go-cursor/internal/service"
	"go-cursor/pkg/database"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-cursor/internal/docs" // This is required for swagger
)

// @title           Go-Cursor API
// @version         1.0
// @description     A User Management API with PostgreSQL and Redis.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
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

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add before routes
	r.Use(cors.Default())

	// Routes
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	// Start server
	r.Run(":8080")
}
