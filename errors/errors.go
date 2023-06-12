package errors

import "fmt"

// CustomError is a custom error type
type CustomError struct {
	Code    int
	Message string
}

// Error returns the error message
func (e *CustomError) Error() string {
	return e.Message
}

// New create a general custom error
func New(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    500,
		Message: fmt.Sprintf(s, v...),
	}
}

// NotFound is returned when a resource is not found
func NotFound(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    404,
		Message: fmt.Sprintf(s, v...),
	}
}

// BadRequest is returned when a request is invalid
func BadRequest(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    400,
		Message: fmt.Sprintf(s, v...),
	}
}

// Unauthorized is returned when a request is unauthorized
func Unauthorized(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    401,
		Message: fmt.Sprintf(s, v...),
	}
}

// Forbidden is returned when a request is forbidden
func Forbidden(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    403,
		Message: fmt.Sprintf(s, v...),
	}
}

// Conflict  is returned when a conflict occurs
func Conflict(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    409,
		Message: fmt.Sprintf(s, v...),
	}
}
