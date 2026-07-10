package errors

import (
	"fmt"
	"net/http"
)

// AppError is the application's standard error type. It carries an HTTP status
// code, a user-facing message, and an optional underlying cause.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // internal cause — never serialised to JSON
}

// Error implements the built-in error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap allows errors.Is / errors.As to traverse the cause chain.
func (e *AppError) Unwrap() error {
	return e.Err
}

// ── Predefined sentinel errors ────────────────────────────────────────────────

// ErrUnauthorized is returned when a request lacks valid authentication credentials.
var ErrUnauthorized = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}

// ErrForbidden is returned when an authenticated user lacks permission for the action.
var ErrForbidden = &AppError{Code: http.StatusForbidden, Message: "forbidden"}

// ErrNotFound is returned when the requested resource does not exist.
var ErrNotFound = &AppError{Code: http.StatusNotFound, Message: "resource not found"}

// ErrConflict is returned when the request conflicts with existing data (e.g. duplicate email).
var ErrConflict = &AppError{Code: http.StatusConflict, Message: "conflict"}

// ErrValidation is returned when request body fails validation.
var ErrValidation = &AppError{Code: http.StatusUnprocessableEntity, Message: "validation error"}

// ErrInternalServer is returned for unexpected server-side failures.
var ErrInternalServer = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}

// ── Constructors ──────────────────────────────────────────────────────────────

// NewAppError creates a new *AppError with a custom HTTP code, message, and optional cause.
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// WrapError wraps any error as a 500 Internal Server Error AppError,
// preserving the original cause for logging while surfacing only the
// provided message to callers.
func WrapError(err error, message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}
