package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestAPIStatusError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIStatusError
		want string
	}{
		{
			name: "with underlying error",
			err:  &APIStatusError{StatusCode: 400, Body: "bad request", Err: errors.New("validation failed")},
			want: "cerebras: API status error 400: Bad Request: validation failed",
		},
		{
			name: "without underlying error",
			err:  &APIStatusError{StatusCode: 404, Body: "not found"},
			want: "cerebras: API status error 404: Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("APIStatusError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIStatusError_IsAPIError(t *testing.T) {
	err := &APIStatusError{StatusCode: 400}
	if !err.IsAPIError() {
		t.Error("APIStatusError.IsAPIError() should return true")
	}
}

func TestAPIStatusError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &APIStatusError{StatusCode: 400, Err: underlying}

	if got := err.Unwrap(); got != underlying {
		t.Errorf("APIStatusError.Unwrap() = %v, want %v", got, underlying)
	}
}

func TestSpecificErrorTypes_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"BadRequestError matches itself", &BadRequestError{}, &BadRequestError{}, true},
		{"BadRequestError doesn't match others", &BadRequestError{}, &NotFoundError{}, false},
		{"AuthenticationError matches itself", &AuthenticationError{}, &AuthenticationError{}, true},
		{"PermissionDeniedError matches itself", &PermissionDeniedError{}, &PermissionDeniedError{}, true},
		{"NotFoundError matches itself", &NotFoundError{}, &NotFoundError{}, true},
		{"UnprocessableEntityError matches itself", &UnprocessableEntityError{}, &UnprocessableEntityError{}, true},
		{"RateLimitError matches itself", &RateLimitError{}, &RateLimitError{}, true},
		{"InternalServerError matches itself", &InternalServerError{}, &InternalServerError{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIConnectionError(t *testing.T) {
	underlying := &url.Error{Op: "Get", URL: "https://api.example.com", Err: errors.New("connection refused")}
	err := &APIConnectionError{Err: underlying}

	t.Run("Error message", func(t *testing.T) {
		want := "cerebras: API connection error: Get \"https://api.example.com\": connection refused"
		if got := err.Error(); got != want {
			t.Errorf("APIConnectionError.Error() = %v, want %v", got, want)
		}
	})

	t.Run("IsAPIError", func(t *testing.T) {
		if !err.IsAPIError() {
			t.Error("APIConnectionError.IsAPIError() should return true")
		}
	})

	t.Run("Is", func(t *testing.T) {
		if !errors.Is(err, &APIConnectionError{}) {
			t.Error("APIConnectionError should match itself")
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		if got := err.Unwrap(); got != underlying {
			t.Errorf("APIConnectionError.Unwrap() = %v, want %v", got, underlying)
		}
	})
}

func TestAPITimeoutError(t *testing.T) {
	err := &APITimeoutError{Err: context.DeadlineExceeded}

	t.Run("Error message", func(t *testing.T) {
		want := "cerebras: API timeout error: context deadline exceeded"
		if got := err.Error(); got != want {
			t.Errorf("APITimeoutError.Error() = %v, want %v", got, want)
		}
	})

	t.Run("IsAPIError", func(t *testing.T) {
		if !err.IsAPIError() {
			t.Error("APITimeoutError.IsAPIError() should return true")
		}
	})

	t.Run("Is", func(t *testing.T) {
		if !errors.Is(err, &APITimeoutError{}) {
			t.Error("APITimeoutError should match itself")
		}
	})
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"RateLimitError", &RateLimitError{}, true},
		{"InternalServerError", &InternalServerError{}, true},
		{"APIConnectionError", &APIConnectionError{Err: errors.New("connection refused")}, true},
		{"APITimeoutError", &APITimeoutError{Err: context.DeadlineExceeded}, true},
		{"BadRequestError", &BadRequestError{}, false},
		{"AuthenticationError", &AuthenticationError{}, false},
		{"PermissionDeniedError", &PermissionDeniedError{}, false},
		{"NotFoundError", &NotFoundError{}, false},
		{"UnprocessableEntityError", &UnprocessableEntityError{}, false},
		{"APIStatusError 408", &APIStatusError{StatusCode: 408}, true},
		{"APIStatusError 409", &APIStatusError{StatusCode: 409}, true},
		{"APIStatusError 429", &APIStatusError{StatusCode: 429}, true},
		{"APIStatusError 500", &APIStatusError{StatusCode: 500}, true},
		{"APIStatusError 502", &APIStatusError{StatusCode: 502}, true},
		{"APIStatusError 503", &APIStatusError{StatusCode: 503}, true},
		{"APIStatusError 400", &APIStatusError{StatusCode: 400}, false},
		{"APIStatusError 401", &APIStatusError{StatusCode: 401}, false},
		{"APIStatusError 403", &APIStatusError{StatusCode: 403}, false},
		{"APIStatusError 404", &APIStatusError{StatusCode: 404}, false},
		{"url.Error", &url.Error{Op: "Get", Err: errors.New("network error")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryableError(tt.err); got != tt.want {
				t.Errorf("IsRetryableError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantType   error
	}{
		{"400 Bad Request", http.StatusBadRequest, "bad request", &BadRequestError{}},
		{"401 Unauthorized", http.StatusUnauthorized, "unauthorized", &AuthenticationError{}},
		{"403 Forbidden", http.StatusForbidden, "forbidden", &PermissionDeniedError{}},
		{"404 Not Found", http.StatusNotFound, "not found", &NotFoundError{}},
		{"422 Unprocessable Entity", http.StatusUnprocessableEntity, "unprocessable", &UnprocessableEntityError{}},
		{"429 Too Many Requests", http.StatusTooManyRequests, "rate limited", &RateLimitError{}},
		{"500 Internal Server Error", http.StatusInternalServerError, "server error", &InternalServerError{}},
		{"502 Bad Gateway", http.StatusBadGateway, "bad gateway", &InternalServerError{}},
		{"503 Service Unavailable", http.StatusServiceUnavailable, "unavailable", &InternalServerError{}},
		{"504 Gateway Timeout", http.StatusGatewayTimeout, "timeout", &InternalServerError{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapStatusCode(tt.statusCode, tt.body)
			// Use type assertion instead of errors.As with interface
			gotType := fmt.Sprintf("%T", err)
			wantType := fmt.Sprintf("%T", tt.wantType)
			if gotType != wantType {
				t.Errorf("MapStatusCode(%d) = %T, want %T", tt.statusCode, err, tt.wantType)
			}
		})
	}
}

func TestNewConnectionError(t *testing.T) {
	underlying := errors.New("network error")
	err := NewConnectionError(underlying)

	if err.Err != underlying {
		t.Errorf("NewConnectionError().Err = %v, want %v", err.Err, underlying)
	}

	if !errors.Is(err, &APIConnectionError{}) {
		t.Error("NewConnectionError should create APIConnectionError")
	}
}

func TestNewTimeoutError(t *testing.T) {
	t.Run("context.DeadlineExceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		<-ctx.Done()

		err := NewTimeoutError(ctx)
		if err.Err != context.DeadlineExceeded {
			t.Errorf("NewTimeoutError().Err = %v, want context.DeadlineExceeded", err.Err)
		}
	})

	t.Run("context.Canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := NewTimeoutError(ctx)
		if err.Err != context.Canceled {
			t.Errorf("NewTimeoutError().Err = %v, want context.Canceled", err.Err)
		}
	})
}
