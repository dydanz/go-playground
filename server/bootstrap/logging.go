package bootstrap

import (
	"go-playground/pkg/logging"
)

// InitializeLogging sets up the logging for the application
func InitializeLogging() {
	logging.InitLogger()
}
