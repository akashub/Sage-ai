// // backend/go/pkg/errors/errors.go
// package errors

// type ValidationError struct {
//     ValidationResult map[string]interface{}
//     Message         string
// }

// func (e *ValidationError) Error() string {
//     return e.Message
// }

// func NewValidationError(result map[string]interface{}, message string) *ValidationError {
//     return &ValidationError{
//         ValidationResult: result,
//         Message:         message,
//     }
// }

// backend/go/pkg/errors/errors.go
package errors

import (
	"fmt"
)

// ValidationError represents an error that occurs during query validation
type ValidationError struct {
	Message          string
	ValidationResult map[string]interface{}
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, result map[string]interface{}) *ValidationError {
	return &ValidationError{
		Message:          message,
		ValidationResult: result,
	}
}

// ExecutionError represents an error that occurs during query execution
type ExecutionError struct {
	Message string
	Query   string
}

// Error returns the error message
func (e *ExecutionError) Error() string {
	return fmt.Sprintf("execution error: %s (query: %s)", e.Message, e.Query)
}

// NewExecutionError creates a new ExecutionError
func NewExecutionError(message, query string) *ExecutionError {
	return &ExecutionError{
		Message: message,
		Query:   query,
	}
}