package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger initializes the zerolog logger with pretty console output
func InitLogger() {
	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create console writer with pretty print
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Initialize global logger
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = Logger
}

// GetLogger returns the configured logger instance
func GetLogger() zerolog.Logger {
	return Logger
}
