package domain

import (
	"errors"
	"fmt"
)

// Common error types
type (
	ValidationError struct {
		Field   string
		Message string
	}

	ResourceNotFoundError struct {
		Resource string
		ID       string
		Message  string
	}

	AuthenticationError struct {
		Message string
	}

	AuthorizationError struct {
		Message string
	}

	ResourceConflictError struct {
		Resource string
		Message  string
	}

	RateLimitError struct {
		Message string
	}

	BusinessLogicError struct {
		Code    string
		Message string
	}

	SystemError struct {
		Op      string
		Err     error
		Message string
	}
)

// Error implementations for each type
func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e ResourceNotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID %s not found: %s", e.Resource, e.ID, e.Message)
	}
	return fmt.Sprintf("%s not found: %s", e.Resource, e.Message)
}

func (e AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

func (e AuthorizationError) Error() string {
	return fmt.Sprintf("authorization error: %s", e.Message)
}

func (e ResourceConflictError) Error() string {
	return fmt.Sprintf("%s conflict: %s", e.Resource, e.Message)
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: %s", e.Message)
}

func (e BusinessLogicError) Error() string {
	return fmt.Sprintf("business logic error [%s]: %s", e.Code, e.Message)
}

func (e SystemError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("system error in %s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("system error in %s: %s", e.Op, e.Message)
}

// Error type checking functions
func IsValidationError(err error) bool {
	var e ValidationError
	return errors.As(err, &e)
}

func IsResourceNotFoundError(err error) bool {
	var e ResourceNotFoundError
	return errors.As(err, &e)
}

func IsAuthenticationError(err error) bool {
	var e AuthenticationError
	return errors.As(err, &e)
}

func IsAuthorizationError(err error) bool {
	var e AuthorizationError
	return errors.As(err, &e)
}

func IsResourceConflictError(err error) bool {
	var e ResourceConflictError
	return errors.As(err, &e)
}

func IsRateLimitError(err error) bool {
	var e RateLimitError
	return errors.As(err, &e)
}

func IsBusinessLogicError(err error) bool {
	var e BusinessLogicError
	return errors.As(err, &e)
}

func IsSystemError(err error) bool {
	var e SystemError
	return errors.As(err, &e)
}

// Helper functions to create errors
func NewValidationError(field, message string) error {
	return ValidationError{Field: field, Message: message}
}

func NewResourceNotFoundError(resource, id, message string) error {
	return ResourceNotFoundError{Resource: resource, ID: id, Message: message}
}

func NewAuthenticationError(message string) error {
	return AuthenticationError{Message: message}
}

func NewAuthorizationError(message string) error {
	return AuthorizationError{Message: message}
}

func NewResourceConflictError(resource, message string) error {
	return ResourceConflictError{Resource: resource, Message: message}
}

func NewRateLimitError(message string) error {
	return RateLimitError{Message: message}
}

func NewBusinessLogicError(code, message string) error {
	return BusinessLogicError{Code: code, Message: message}
}

func NewSystemError(op string, err error, message string) error {
	return SystemError{Op: op, Err: err, Message: message}
}
