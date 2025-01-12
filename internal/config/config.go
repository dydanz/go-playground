package config

import (
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

type Config struct {
	// PostgreSQL settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

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
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "go_cursor"),

		// Redis settings
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "your_redis_password_here"),

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
