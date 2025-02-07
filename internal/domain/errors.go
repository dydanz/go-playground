package domain

// Custom error types for authentication
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type AuthenticationError struct {
	Message string
}

func (e AuthenticationError) Error() string {
	return e.Message
}

type ResourceNotFoundError struct {
	Resource string
	Message  string
}

func (e ResourceNotFoundError) Error() string {
	return e.Message
}

type ResourceConflictError struct {
	Resource string
	Message  string
}

func (e ResourceConflictError) Error() string {
	return e.Message
}

type InvalidInputError struct {
	Message string
}

func (e InvalidInputError) Error() string {
	return e.Message
}

// Error checking helper functions
func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

func IsAuthenticationError(err error) bool {
	_, ok := err.(AuthenticationError)
	return ok
}

func IsResourceNotFoundError(err error) bool {
	_, ok := err.(ResourceNotFoundError)
	return ok
}

func IsResourceConflictError(err error) bool {
	_, ok := err.(ResourceConflictError)
	return ok
}

func IsInvalidInputError(err error) bool {
	_, ok := err.(InvalidInputError)
	return ok
}
