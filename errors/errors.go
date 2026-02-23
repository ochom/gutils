// Package errors provides HTTP-aware error types with status codes.
//
// This package defines custom error types that include HTTP status codes,
// making it easy to return appropriate HTTP responses in web applications.
// All error constructors support printf-style formatting.
//
// Example usage:
//
//	func GetUser(id string) (*User, error) {
//		user, err := db.FindUser(id)
//		if err != nil {
//			return nil, errors.NotFound("user with id %s not found", id)
//		}
//		return user, nil
//	}
//
//	func CreateUser(data UserData) error {
//		if data.Email == "" {
//			return errors.BadRequest("email is required")
//		}
//		if existingUser != nil {
//			return errors.Conflict("user with email %s already exists", data.Email)
//		}
//		return nil
//	}
//
//	// In HTTP handler
//	func Handler(w http.ResponseWriter, r *http.Request) {
//		user, err := GetUser(id)
//		if err != nil {
//			if customErr, ok := err.(*errors.CustomError); ok {
//				http.Error(w, customErr.Message, customErr.Code)
//				return
//			}
//		}
//	}
package errors

import "fmt"

// CustomError represents an error with an associated HTTP status code.
// It implements the error interface and can be used anywhere a standard error is expected.
type CustomError struct {
	// Code is the HTTP status code associated with this error
	Code int
	// Message is the human-readable error message
	Message string
}

// Error implements the error interface, returning the error message.
func (e *CustomError) Error() string {
	return e.Message
}

// New creates a general server error (HTTP 500) with the given formatted message.
//
// Example:
//
//	return errors.New("database connection failed: %v", err)
//	return errors.New("unexpected error occurred")
func New(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    500,
		Message: fmt.Sprintf(s, v...),
	}
}

// NotFound creates a 404 Not Found error.
// Use when a requested resource does not exist.
//
// Example:
//
//	return errors.NotFound("user %s not found", userID)
//	return errors.NotFound("article with slug %s does not exist", slug)
func NotFound(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    404,
		Message: fmt.Sprintf(s, v...),
	}
}

// BadRequest creates a 400 Bad Request error.
// Use when the client's request is malformed or contains invalid data.
//
// Example:
//
//	return errors.BadRequest("invalid email format: %s", email)
//	return errors.BadRequest("age must be a positive number")
func BadRequest(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    400,
		Message: fmt.Sprintf(s, v...),
	}
}

// Unauthorized creates a 401 Unauthorized error.
// Use when authentication is required but not provided or invalid.
//
// Example:
//
//	return errors.Unauthorized("invalid or expired token")
//	return errors.Unauthorized("authentication required")
func Unauthorized(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    401,
		Message: fmt.Sprintf(s, v...),
	}
}

// Forbidden creates a 403 Forbidden error.
// Use when the user is authenticated but lacks permission to access the resource.
//
// Example:
//
//	return errors.Forbidden("you do not have access to this resource")
//	return errors.Forbidden("admin privileges required")
func Forbidden(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    403,
		Message: fmt.Sprintf(s, v...),
	}
}

// Conflict creates a 409 Conflict error.
// Use when the request conflicts with the current state of the resource.
//
// Example:
//
//	return errors.Conflict("email %s is already registered", email)
//	return errors.Conflict("resource version conflict, please refresh and retry")
func Conflict(s string, v ...any) *CustomError {
	return &CustomError{
		Code:    409,
		Message: fmt.Sprintf(s, v...),
	}
}
