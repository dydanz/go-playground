package bootstrap

import (
	"go-playground/server/config"
	"go-playground/pkg/database"
	"log"

	"github.com/go-redis/redis/v8"
)

// InitializeDatabase initializes both primary and replication database connections
func InitializeDatabase(cfg *config.Config) *config.DbConnection {
	// Initialize PostgreSQL Primary
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize PostgreSQL Replication
	dbReplication, err := database.NewPostgresReplicationConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to replication database: %v", err)
	}

	return &config.DbConnection{
		RW: db,
		RR: dbReplication,
	}
}

// InitializeRedis initializes Redis connection
func InitializeRedis(cfg *config.Config) *redis.Client {
	return database.NewRedisConnection(cfg)
}
