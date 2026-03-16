// Package errors provides custom error types for the Cerebras Cloud SDK.
//
// The error hierarchy mirrors the Python SDK for easy migration:
//
//	APIError (base)
//	├── APIStatusError (HTTP status errors)
//	│   ├── BadRequestError (400)
//	│   ├── AuthenticationError (401)
//	│   ├── PermissionDeniedError (403)
//	│   ├── NotFoundError (404)
//	│   ├── UnprocessableEntityError (422)
//	│   ├── RateLimitError (429)
//	│   └── InternalServerError (>=500)
//	└── APIConnectionError (network errors)
//	    └── APITimeoutError (context.DeadlineExceeded)
//
// All error types support errors.Is() and errors.As() for programmatic handling.
// Use IsRetryableError() to check if an error should trigger automatic retry.
package errors

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	stderrors "errors"
)

// APIError is the base interface for all Cerebras API errors.
// All error types in this package implement this interface.
type APIError interface {
	error
	// IsAPIError marks this as an API error (for type assertions)
	IsAPIError() bool
}

// APIStatusError represents an error response with an HTTP status code.
type APIStatusError struct {
	StatusCode int
	Body       string
	Err        error
}

// Error implements the error interface.
func (e *APIStatusError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cerebras: API status error %d: %s: %s", e.StatusCode, http.StatusText(e.StatusCode), e.Err.Error())
	}
	return fmt.Sprintf("cerebras: API status error %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// Unwrap implements errors.Unwrap() for errors.Is/As support.
func (e *APIStatusError) Unwrap() error {
	return e.Err
}

// IsAPIError implements the APIError interface.
func (e *APIStatusError) IsAPIError() bool {
	return true
}

// BadRequestError represents a 400 Bad Request response.
type BadRequestError struct {
	APIStatusError
}

// Is implements errors.Is() for type checking.
func (e *BadRequestError) Is(target error) bool {
	_, ok := target.(*BadRequestError)
	return ok
}

// AuthenticationError represents a 401 Unauthorized response.
type AuthenticationError struct {
	APIStatusError
}

// Is implements errors.Is() for type checking.
func (e *AuthenticationError) Is(target error) bool {
	_, ok := target.(*AuthenticationError)
	return ok
}

// PermissionDeniedError represents a 403 Forbidden response.
type PermissionDeniedError struct {
	APIStatusError
}

// Is implements errors.Is() for type checking.
func (e *PermissionDeniedError) Is(target error) bool {
	_, ok := target.(*PermissionDeniedError)
	return ok
}

// NotFoundError represents a 404 Not Found response.
type NotFoundError struct {
	APIStatusError
}

// Is implements errors.Is() for type checking.
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

// UnprocessableEntityError represents a 422 Unprocessable Entity response.
type UnprocessableEntityError struct {
	APIStatusError
}

// Is implements errors.Is() for type checking.
func (e *UnprocessableEntityError) Is(target error) bool {
	_, ok := target.(*UnprocessableEntityError)
	return ok
}

// RateLimitError represents a 429 Too Many Requests response.
type RateLimitError struct {
	APIStatusError
	RetryAfter int // seconds to wait before retrying (if provided)
}

// Is implements errors.Is() for type checking.
func (e *RateLimitError) Is(target error) bool {
	_, ok := target.(*RateLimitError)
	return ok
}

// InternalServerError represents a 5xx Internal Server Error response.
type InternalServerError struct {
	APIStatusError
}

// Is implements errors.Is() for type checking.
func (e *InternalServerError) Is(target error) bool {
	_, ok := target.(*InternalServerError)
	return ok
}

// APIConnectionError represents a network-level error.
type APIConnectionError struct {
	Err error
}

// Error implements the error interface.
func (e *APIConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cerebras: API connection error: %s", e.Err.Error())
	}
	return "cerebras: API connection error"
}

// Unwrap implements errors.Unwrap() for errors.Is/As support.
func (e *APIConnectionError) Unwrap() error {
	return e.Err
}

// IsAPIError implements the APIError interface.
func (e *APIConnectionError) IsAPIError() bool {
	return true
}

// Is implements errors.Is() for type checking.
func (e *APIConnectionError) Is(target error) bool {
	_, ok := target.(*APIConnectionError)
	return ok
}

// APITimeoutError represents a request timeout.
type APITimeoutError struct {
	Err error
}

// Error implements the error interface.
func (e *APITimeoutError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cerebras: API timeout error: %s", e.Err.Error())
	}
	return "cerebras: API timeout error"
}

// Unwrap implements errors.Unwrap() for errors.Is/As support.
func (e *APITimeoutError) Unwrap() error {
	return e.Err
}

// IsAPIError implements the APIError interface.
func (e *APITimeoutError) IsAPIError() bool {
	return true
}

// Is implements errors.Is() for type checking.
func (e *APITimeoutError) Is(target error) bool {
	_, ok := target.(*APITimeoutError)
	return ok
}

// IsRetryableError returns true if the error is eligible for automatic retry.
//
// Retryable errors include:
//   - RateLimitError (429)
//   - InternalServerError (5xx)
//   - APIConnectionError (network errors)
//   - APITimeoutError (timeouts)
//
// Non-retryable errors include:
//   - BadRequestError (400)
//   - AuthenticationError (401)
//   - PermissionDeniedError (403)
//   - NotFoundError (404)
//   - UnprocessableEntityError (422)
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific retryable error types
	var rateLimitErr *RateLimitError
	if stderrors.As(err, &rateLimitErr) {
		return true
	}

	var internalErr *InternalServerError
	if stderrors.As(err, &internalErr) {
		return true
	}

	var connectionErr *APIConnectionError
	if stderrors.As(err, &connectionErr) {
		return true
	}

	var timeoutErr *APITimeoutError
	if stderrors.As(err, &timeoutErr) {
		return true
	}

	// Check for generic APIStatusError with retryable status codes
	var statusErr *APIStatusError
	if stderrors.As(err, &statusErr) {
		// Retry on 408, 409, 429, and 5xx
		if statusErr.StatusCode == 408 || statusErr.StatusCode == 409 || statusErr.StatusCode == 429 {
			return true
		}
		if statusErr.StatusCode >= 500 && statusErr.StatusCode < 600 {
			return true
		}
	}

	// Check for url.Error (network errors)
	var urlErr *url.Error
	if stderrors.As(err, &urlErr) {
		return true
	}

	return false
}

// MapStatusCode maps an HTTP status code to the appropriate error type.
func MapStatusCode(statusCode int, body string) error {
	switch statusCode {
	case http.StatusBadRequest: // 400
		return &BadRequestError{
			APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
		}
	case http.StatusUnauthorized: // 401
		return &AuthenticationError{
			APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
		}
	case http.StatusForbidden: // 403
		return &PermissionDeniedError{
			APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
		}
	case http.StatusNotFound: // 404
		return &NotFoundError{
			APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
		}
	case http.StatusUnprocessableEntity: // 422
		return &UnprocessableEntityError{
			APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
		}
	case http.StatusTooManyRequests: // 429
		return &RateLimitError{
			APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
		}
	default:
		if statusCode >= 500 && statusCode < 600 {
			return &InternalServerError{
				APIStatusError: APIStatusError{StatusCode: statusCode, Body: body},
			}
		}
		// For other status codes, return a generic APIStatusError
		return &APIStatusError{StatusCode: statusCode, Body: body}
	}
}

// NewConnectionError creates a new APIConnectionError from an underlying error.
func NewConnectionError(err error) *APIConnectionError {
	return &APIConnectionError{Err: err}
}

// NewTimeoutError creates a new APITimeoutError from a context error.
func NewTimeoutError(ctx context.Context) *APITimeoutError {
	if ctx.Err() == context.DeadlineExceeded {
		return &APITimeoutError{Err: context.DeadlineExceeded}
	}
	return &APITimeoutError{Err: ctx.Err()}
}
