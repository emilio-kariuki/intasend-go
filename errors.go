package intasend

import (
	"errors"
	"fmt"
)

// Sentinel errors for common error conditions.
var (
	ErrMissingPublishableKey = errors.New("intasend: publishable key is required")
	ErrMissingSecretKey      = errors.New("intasend: secret key is required")
	ErrInvalidEnvironment    = errors.New("intasend: could not determine environment from keys")
	ErrNoKeysProvided        = errors.New("intasend: at least one API key must be provided")
)

// APIError represents an error returned by the IntaSend API.
type APIError struct {
	// HTTPStatusCode is the HTTP status code of the response.
	HTTPStatusCode int `json:"-"`

	// Code is the IntaSend error code, if provided.
	Code string `json:"code,omitempty"`

	// Message is the human-readable error message.
	Message string `json:"message,omitempty"`

	// Detail provides additional error details.
	Detail string `json:"detail,omitempty"`

	// Errors contains field-level validation errors.
	Errors map[string][]string `json:"errors,omitempty"`

	// RequestID is the unique request identifier for debugging.
	RequestID string `json:"request_id,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("intasend: API error (HTTP %d): %s", e.HTTPStatusCode, e.Message)
	}
	if e.Detail != "" {
		return fmt.Sprintf("intasend: API error (HTTP %d): %s", e.HTTPStatusCode, e.Detail)
	}
	return fmt.Sprintf("intasend: API error (HTTP %d)", e.HTTPStatusCode)
}

// IsNotFound returns true if the error indicates a resource was not found.
func (e *APIError) IsNotFound() bool {
	return e.HTTPStatusCode == 404
}

// IsAuthenticationError returns true if the error is an authentication failure.
func (e *APIError) IsAuthenticationError() bool {
	return e.HTTPStatusCode == 401 || e.HTTPStatusCode == 403
}

// IsValidationError returns true if the error is a validation error.
func (e *APIError) IsValidationError() bool {
	return e.HTTPStatusCode == 400 && len(e.Errors) > 0
}

// IsRateLimited returns true if the request was rate limited.
func (e *APIError) IsRateLimited() bool {
	return e.HTTPStatusCode == 429
}

// NetworkError represents a network-level error.
type NetworkError struct {
	Err     error
	Message string
}

// Error implements the error interface.
func (e *NetworkError) Error() string {
	return fmt.Sprintf("intasend: network error: %s: %v", e.Message, e.Err)
}

// Unwrap returns the underlying error.
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// IsAPIError checks if an error is an IntaSend API error.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsNetworkError checks if an error is a network error.
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

// AsAPIError attempts to extract an APIError from the given error.
// Returns nil if the error is not an APIError.
func AsAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}
