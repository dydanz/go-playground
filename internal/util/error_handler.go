package util

import (
	"go-playground/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
	Type  string `json:"type"`
}

// HandleError is a helper function to handle different types of errors and return appropriate HTTP status codes
func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case domain.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
	case domain.ResourceNotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": e.Message})
	case domain.ResourceConflictError:
		c.JSON(http.StatusConflict, gin.H{"error": e.Message})
	case domain.AuthenticationError:
		c.JSON(http.StatusUnauthorized, gin.H{"error": e.Message})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GetErrorStatusCode returns the appropriate HTTP status code for an error without sending a response
func GetErrorStatusCode(err error) int {
	switch {
	case domain.IsValidationError(err):
		return http.StatusBadRequest
	case domain.IsAuthenticationError(err):
		return http.StatusUnauthorized
	case domain.IsResourceNotFoundError(err):
		return http.StatusNotFound
	case domain.IsResourceConflictError(err):
		return http.StatusConflict
	case domain.IsInvalidInputError(err):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// GetErrorResponse returns a formatted error response without sending it
func GetErrorResponse(err error) ErrorResponse {
	switch {
	case domain.IsValidationError(err):
		return ErrorResponse{
			Error: err.Error(),
			Type:  "validation_error",
		}
	case domain.IsAuthenticationError(err):
		return ErrorResponse{
			Error: err.Error(),
			Type:  "authentication_error",
		}
	case domain.IsResourceNotFoundError(err):
		return ErrorResponse{
			Error: err.Error(),
			Type:  "not_found_error",
		}
	case domain.IsResourceConflictError(err):
		return ErrorResponse{
			Error: err.Error(),
			Type:  "conflict_error",
		}
	case domain.IsInvalidInputError(err):
		return ErrorResponse{
			Error: err.Error(),
			Type:  "invalid_input_error",
		}
	default:
		return ErrorResponse{
			Error: "An internal server error occurred",
			Type:  "internal_server_error",
		}
	}
}
