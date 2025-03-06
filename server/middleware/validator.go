package middleware

import (
	"fmt"
	"go-playground/pkg/logging"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// RequestValidator interface defines methods for request validation
type RequestValidator interface {
	ValidateRequest(c *gin.Context) error
	GetRules() map[string][]string
}

// BaseValidator provides common validation functionality
type BaseValidator struct {
	validate *validator.Validate
	logger   zerolog.Logger
}

// NewBaseValidator creates a new BaseValidator instance
func NewBaseValidator() *BaseValidator {
	return &BaseValidator{
		validate: validator.New(),
		logger:   logging.GetLogger(),
	}
}

// ValidateStruct validates a struct using validator tags
func (v *BaseValidator) ValidateStruct(obj interface{}) error {
	if err := v.validate.Struct(obj); err != nil {
		var errors ValidationErrors
		for _, err := range err.(validator.ValidationErrors) {
			v.logger.Debug().
				Str("field", err.Field()).
				Str("tag", err.Tag()).
				Str("value", fmt.Sprintf("%v", err.Value())).
				Msg("Validation error")

			errors = append(errors, ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: getErrorMsg(err),
			})
		}
		return errors
	}
	return nil
}

// RequestValidationMiddleware is a middleware that validates requests
func RequestValidationMiddleware(validator RequestValidator) gin.HandlerFunc {
	logger := logging.GetLogger()

	return func(c *gin.Context) {
		if err := validator.ValidateRequest(c); err != nil {
			var validationErrors ValidationErrors
			if errors, ok := err.(ValidationErrors); ok {
				validationErrors = errors
			} else {
				validationErrors = ValidationErrors{{
					Field:   "request",
					Message: err.Error(),
				}}
			}

			logger.Error().
				Str("method", c.Request.Method).
				Str("url", c.Request.URL.RequestURI()).
				Interface("errors", validationErrors).
				Msg("Request validation failed")

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationErrors,
			})
			c.Abort()
			return
		}

		logger.Debug().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.RequestURI()).
			Msg("Request validation successful")

		c.Next()
	}
}

// getErrorMsg returns a human-readable error message for validation errors
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must not be longer than %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	case "uuid":
		return "Must be a valid UUID"
	case "datetime":
		return "Must be a valid datetime"
	case "gt":
		return fmt.Sprintf("Must be greater than %s", err.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", err.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", err.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", err.Param())
	case "url":
		return "Must be a valid URL"
	case "alpha":
		return "Must contain only letters"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "numeric":
		return "Must be numeric"
	case "boolean":
		return "Must be a boolean value"
	default:
		return fmt.Sprintf("Failed on %s validation", err.Tag())
	}
}

// Example validator for JSON payloads
type JSONValidator struct {
	*BaseValidator
	rules map[string][]string
}

// NewJSONValidator creates a new JSONValidator instance
func NewJSONValidator(rules map[string][]string) *JSONValidator {
	return &JSONValidator{
		BaseValidator: NewBaseValidator(),
		rules:         rules,
	}
}

func (v *JSONValidator) ValidateRequest(c *gin.Context) error {
	// Get content type
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("content-type must be application/json")
	}

	// Read and validate the request body
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	// Validate required fields and their types
	var errors ValidationErrors
	for field, rules := range v.rules {
		value, exists := body[field]

		// Skip validation for optional fields that are not present
		if !exists && !containsRule(rules, "required") {
			continue
		}

		for _, rule := range rules {
			// Skip empty optional fields
			if !exists && rule != "required" {
				continue
			}

			switch rule {
			case "required":
				if !exists || value == nil {
					errors = append(errors, ValidationError{
						Field:   field,
						Message: "This field is required",
					})
				}
			case "string":
				if exists && value != nil && reflect.TypeOf(value).Kind() != reflect.String {
					errors = append(errors, ValidationError{
						Field:   field,
						Message: "Must be a string",
					})
				}
			case "number":
				if exists && value != nil {
					switch reflect.TypeOf(value).Kind() {
					case reflect.Float64, reflect.Float32, reflect.Int, reflect.Int64:
						// Valid number types in JSON
					default:
						errors = append(errors, ValidationError{
							Field:   field,
							Message: "Must be a number",
						})
					}
				}
			case "boolean":
				if exists && value != nil && reflect.TypeOf(value).Kind() != reflect.Bool {
					errors = append(errors, ValidationError{
						Field:   field,
						Message: "Must be a boolean",
					})
				}
			case "array":
				if exists && value != nil && reflect.TypeOf(value).Kind() != reflect.Slice {
					errors = append(errors, ValidationError{
						Field:   field,
						Message: "Must be an array",
					})
				}
			case "object":
				if exists && value != nil && reflect.TypeOf(value).Kind() != reflect.Map {
					errors = append(errors, ValidationError{
						Field:   field,
						Message: "Must be an object",
					})
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// Helper function to check if a rule exists in the rules slice
func containsRule(rules []string, rule string) bool {
	for _, r := range rules {
		if r == rule {
			return true
		}
	}
	return false
}

func (v *JSONValidator) GetRules() map[string][]string {
	return v.rules
}
