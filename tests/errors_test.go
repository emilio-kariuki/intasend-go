package tests

import (
	"errors"
	"fmt"
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

func TestAPIError_ErrorMessage(t *testing.T) {
	e := &intasend.APIError{HTTPStatusCode: 400, Message: "bad request"}
	want := "intasend: API error (HTTP 400): bad request"
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}

func TestAPIError_ErrorDetail(t *testing.T) {
	e := &intasend.APIError{HTTPStatusCode: 500, Detail: "internal error"}
	want := "intasend: API error (HTTP 500): internal error"
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}

func TestAPIError_ErrorFallback(t *testing.T) {
	e := &intasend.APIError{HTTPStatusCode: 503}
	want := "intasend: API error (HTTP 503)"
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}

func TestAPIError_MessagePrecedence(t *testing.T) {
	e := &intasend.APIError{HTTPStatusCode: 400, Message: "msg", Detail: "detail"}
	// Message should take precedence over Detail
	if e.Error() != "intasend: API error (HTTP 400): msg" {
		t.Errorf("message should take precedence, got %q", e.Error())
	}
}

func TestAPIError_IsNotFound(t *testing.T) {
	tests := []struct {
		status int
		want   bool
	}{
		{404, true},
		{200, false},
		{400, false},
	}
	for _, tt := range tests {
		e := &intasend.APIError{HTTPStatusCode: tt.status}
		if got := e.IsNotFound(); got != tt.want {
			t.Errorf("IsNotFound() with status %d = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestAPIError_IsAuthenticationError(t *testing.T) {
	tests := []struct {
		status int
		want   bool
	}{
		{401, true},
		{403, true},
		{200, false},
		{400, false},
	}
	for _, tt := range tests {
		e := &intasend.APIError{HTTPStatusCode: tt.status}
		if got := e.IsAuthenticationError(); got != tt.want {
			t.Errorf("IsAuthenticationError() with status %d = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestAPIError_IsValidationError(t *testing.T) {
	e := &intasend.APIError{
		HTTPStatusCode: 400,
		Errors:         map[string][]string{"phone": {"required"}},
	}
	if !e.IsValidationError() {
		t.Error("expected IsValidationError() to be true")
	}

	// 400 without Errors field
	e2 := &intasend.APIError{HTTPStatusCode: 400}
	if e2.IsValidationError() {
		t.Error("expected IsValidationError() to be false without Errors")
	}

	// Non-400 with Errors field
	e3 := &intasend.APIError{
		HTTPStatusCode: 500,
		Errors:         map[string][]string{"field": {"err"}},
	}
	if e3.IsValidationError() {
		t.Error("expected IsValidationError() to be false for non-400 status")
	}
}

func TestAPIError_IsRateLimited(t *testing.T) {
	tests := []struct {
		status int
		want   bool
	}{
		{429, true},
		{200, false},
		{400, false},
	}
	for _, tt := range tests {
		e := &intasend.APIError{HTTPStatusCode: tt.status}
		if got := e.IsRateLimited(); got != tt.want {
			t.Errorf("IsRateLimited() with status %d = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestNetworkError_Error(t *testing.T) {
	inner := fmt.Errorf("connection refused")
	e := &intasend.NetworkError{Err: inner, Message: "request failed"}
	want := "intasend: network error: request failed: connection refused"
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}

func TestNetworkError_Unwrap(t *testing.T) {
	inner := fmt.Errorf("connection refused")
	e := &intasend.NetworkError{Err: inner, Message: "request failed"}
	if e.Unwrap() != inner {
		t.Error("Unwrap() should return inner error")
	}
}

func TestIsAPIError(t *testing.T) {
	apiErr := &intasend.APIError{HTTPStatusCode: 400, Message: "bad"}
	if !intasend.IsAPIError(apiErr) {
		t.Error("expected IsAPIError to be true for APIError")
	}

	wrapped := fmt.Errorf("wrapped: %w", apiErr)
	if !intasend.IsAPIError(wrapped) {
		t.Error("expected IsAPIError to be true for wrapped APIError")
	}

	if intasend.IsAPIError(fmt.Errorf("plain error")) {
		t.Error("expected IsAPIError to be false for plain error")
	}

	if intasend.IsAPIError(nil) {
		t.Error("expected IsAPIError to be false for nil")
	}
}

func TestIsNetworkError(t *testing.T) {
	netErr := &intasend.NetworkError{Err: fmt.Errorf("timeout"), Message: "failed"}
	if !intasend.IsNetworkError(netErr) {
		t.Error("expected IsNetworkError to be true for NetworkError")
	}

	wrapped := fmt.Errorf("wrapped: %w", netErr)
	if !intasend.IsNetworkError(wrapped) {
		t.Error("expected IsNetworkError to be true for wrapped NetworkError")
	}

	if intasend.IsNetworkError(fmt.Errorf("plain error")) {
		t.Error("expected IsNetworkError to be false for plain error")
	}
}

func TestAsAPIError(t *testing.T) {
	apiErr := &intasend.APIError{HTTPStatusCode: 401, Message: "unauthorized"}
	got := intasend.AsAPIError(apiErr)
	if got == nil {
		t.Fatal("expected non-nil APIError")
	}
	if got.HTTPStatusCode != 401 {
		t.Errorf("expected status 401, got %d", got.HTTPStatusCode)
	}

	wrapped := fmt.Errorf("wrapped: %w", apiErr)
	got = intasend.AsAPIError(wrapped)
	if got == nil {
		t.Fatal("expected non-nil APIError from wrapped error")
	}
	if got.Message != "unauthorized" {
		t.Errorf("expected message 'unauthorized', got %q", got.Message)
	}

	if intasend.AsAPIError(fmt.Errorf("plain error")) != nil {
		t.Error("expected nil for plain error")
	}

	if intasend.AsAPIError(nil) != nil {
		t.Error("expected nil for nil error")
	}
}

func TestNetworkError_ImplementsError(t *testing.T) {
	var err error = &intasend.NetworkError{Err: fmt.Errorf("test"), Message: "msg"}
	_ = err.Error()
}

func TestAPIError_ImplementsError(t *testing.T) {
	var err error = &intasend.APIError{HTTPStatusCode: 500}
	_ = err.Error()
}

func TestNetworkError_ErrorsIs(t *testing.T) {
	sentinel := fmt.Errorf("sentinel")
	netErr := &intasend.NetworkError{Err: sentinel, Message: "failed"}
	if !errors.Is(netErr, sentinel) {
		t.Error("expected errors.Is to find sentinel through NetworkError")
	}
}
