package database

import (
	"database/sql"
	"fmt"
	"strconv"
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

	maxOpenConns, _ := strconv.Atoi(cfg.DBSetMaxOpenConn)
	maxIdleConns, _ := strconv.Atoi(cfg.DBSetMaxIdleConn)
	maxLifeTime, _ := strconv.Atoi(cfg.DBSetMaxLifeTimeConn)
	maxIdleTime, _ := strconv.Atoi(cfg.DBSetMaxIdleTimeConn)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)

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
	maxOpenConns, _ := strconv.Atoi(cfg.DBReplicationSetMaxOpenConn)
	maxIdleConns, _ := strconv.Atoi(cfg.DBReplicationSetMaxIdleConn)
	maxLifeTime, _ := strconv.Atoi(cfg.DBReplicationSetMaxLifeTimeConn)
	maxIdleTime, _ := strconv.Atoi(cfg.DBReplicationSetMaxIdleTimeConn)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping replication DB: %w", err)
	}

	return db, nil
}
