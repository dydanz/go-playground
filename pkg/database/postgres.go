package database

import (
	"database/sql"
	"fmt"
	"time"

	"go-playground/server/config"

	_ "github.com/lib/pq"
)

// NewPostgresConnection initializes a PostgreSQL connection with pooling & timeout settings
func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)             // Max open connections
	db.SetMaxIdleConns(10)             // Max idle connections
	db.SetConnMaxLifetime(5 * time.Minute)  // Max connection lifetime
	db.SetConnMaxIdleTime(2 * time.Minute)  // Max idle time before closing

	// Check the database connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}

// NewPostgresReplicationConnection initializes a replication DB connection with pooling
func NewPostgresReplicationConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBReplicationPort, cfg.DBReplicationUser, cfg.DBReplicationPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open replication DB: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping replication DB: %w", err)
	}

	return db, nil
}
