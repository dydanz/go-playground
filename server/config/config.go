package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
	LoginAttemptResetPeriod time.Duration // Duration after which login attempts are reset
	MaxLoginAttempts        int           // Maximum number of failed attempts before locking
	LockDuration            time.Duration // How long to lock the account after max attempts
}

type DbConnection struct {
	RW *sql.DB
	RR *sql.DB
}

type Config struct {
	// PostgreSQL Primary settings
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSetMaxOpenConn     string
	DBSetMaxIdleConn     string
	DBSetMaxLifeTimeConn string
	DBSetMaxIdleTimeConn string

	// PostgreSQL Replication settings
	DBReplicationUser     string
	DBReplicationPassword string
	DBReplicationPort     string

	DBReplicationSetMaxOpenConn     string
	DBReplicationSetMaxIdleConn     string
	DBReplicationSetMaxLifeTimeConn string
	DBReplicationSetMaxIdleTimeConn string

	// Redis settings
	RedisHost     string
	RedisPort     string
	RedisPassword string

	Auth AuthConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		// PostgreSQL settings
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "go_cursor"),
		DBSetMaxOpenConn:     getEnv("DB_SET_MAX_OPEN_CONN", "10"),
		DBSetMaxIdleConn:     getEnv("DB_SET_MAX_IDLE_CONN", "10"),
		DBSetMaxLifeTimeConn: getEnv("DB_SET_MAX_LIFETIME", "5"),
		DBSetMaxIdleTimeConn: getEnv("DB_SET_MAX_DLE_TIMEOUT", "2"),

		// PostgreSQL Replication settings
		DBReplicationUser:               getEnv("DB_REPLICATION_USER", "replicator"),
		DBReplicationPassword:           getEnv("DB_REPLICATION_PASSWORD", "replicator_password"),
		DBReplicationPort:               getEnv("DB_REPLICATION_PORT", "5433"),
		DBReplicationSetMaxOpenConn:     getEnv("DB_REPLICATION_SET_MAX_OPEN_CONN", "10"),
		DBReplicationSetMaxIdleConn:     getEnv("DB_REPLICATION_SET_MAX_IDLE_CONN", "10"),
		DBReplicationSetMaxLifeTimeConn: getEnv("DB_REPLICATION_SET_MAX_LIFETIME", "5"),
		DBReplicationSetMaxIdleTimeConn: getEnv("DB_REPLICATION_SET_MAX_DLE_TIMEOUT", "2"),

		// Redis settings
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "redis123"),

		Auth: AuthConfig{
			LoginAttemptResetPeriod: 24 * time.Hour,   // Reset attempts after 24 hours
			MaxLoginAttempts:        5,                // Lock after 5 failed attempts
			LockDuration:            30 * time.Minute, // Lock for 30 minutes
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
