package util

import (
	"go-playground/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// HandleError is a helper function to handle different types of errors and return appropriate HTTP status codes
func HandleError(c *gin.Context, err error) {
	var response ErrorResponse
	var statusCode int

	switch e := err.(type) {
	case domain.ValidationError:
		statusCode = http.StatusBadRequest
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    "VALIDATION_ERROR",
			Details: map[string]string{
				"field": e.Field,
			},
		}

	case domain.ResourceNotFoundError:
		statusCode = http.StatusNotFound
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    "NOT_FOUND",
			Details: map[string]string{
				"resource": e.Resource,
				"id":       e.ID,
			},
		}

	case domain.AuthenticationError:
		statusCode = http.StatusUnauthorized
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    "UNAUTHORIZED",
		}

	case domain.AuthorizationError:
		statusCode = http.StatusForbidden
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    "FORBIDDEN",
		}

	case domain.ResourceConflictError:
		statusCode = http.StatusConflict
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    "CONFLICT",
			Details: map[string]string{
				"resource": e.Resource,
			},
		}

	case domain.RateLimitError:
		statusCode = http.StatusTooManyRequests
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    "RATE_LIMIT_EXCEEDED",
		}

	case domain.BusinessLogicError:
		// Business logic errors are returned as 200 with error status
		statusCode = http.StatusOK
		response = ErrorResponse{
			Status:  "error",
			Message: e.Error(),
			Code:    e.Code,
		}

	case domain.SystemError:
		statusCode = http.StatusInternalServerError
		// Log the full error with stack trace for system errors
		LogError(e)
		response = ErrorResponse{
			Status:  "error",
			Message: "An internal server error occurred",
			Code:    "INTERNAL_SERVER_ERROR",
		}

	default:
		statusCode = http.StatusInternalServerError
		// Log unexpected errors
		LogError(err)
		response = ErrorResponse{
			Status:  "error",
			Message: "An internal server error occurred",
			Code:    "INTERNAL_SERVER_ERROR",
		}
	}

	c.JSON(statusCode, response)
}

// LogError logs error details for monitoring and debugging
func LogError(err error) {
	// TODO: Implement proper logging with your preferred logging library
	// Example using standard log package:
	// log.Printf("[ERROR] %v", err)
	//
	// For production, consider using structured logging with:
	// - Timestamp
	// - Error message
	// - Stack trace
	// - Request ID
	// - User ID (if available)
	// - Additional context
}

// EmptyResponse returns a 200 OK with empty data when no results are found
func EmptyResponse(c *gin.Context) {
	c.JSON(http.StatusOK, ErrorResponse{
		Status:  "success",
		Message: "No results found",
		Details: []interface{}{},
	})
}
