// backend/go/pkg/errors/errors.go
package errors

type ValidationError struct {
    ValidationResult map[string]interface{}
    Message         string
}

func (e *ValidationError) Error() string {
    return e.Message
}

func NewValidationError(result map[string]interface{}, message string) *ValidationError {
    return &ValidationError{
        ValidationResult: result,
        Message:         message,
    }
}