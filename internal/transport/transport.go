// Package transport provides an HTTP transport layer with retry logic for the Cerebras Cloud SDK.
//
// The Transport struct wraps net/http.Client and adds:
//   - Automatic retry with exponential backoff and jitter
//   - Request/response JSON marshaling/unmarshaling
//   - Error mapping to custom error types
//   - Context-aware timeout handling
//   - TCP connection warming support
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	cerebraserrors "github/ldaidone/cerebras-cloud-sdk-go/internal/errors"
)

// Default configuration values.
const (
	DefaultMaxRetries     = 3
	DefaultTimeout        = 60 * time.Second
	DefaultBaseBackoff    = 100 * time.Millisecond
	DefaultMaxBackoff     = 2 * time.Second
	DefaultJitterPercent  = 0.25
	TCPWarmingTimeout     = 5 * time.Second
	TCPWarmingEndpoint    = "/v1/tcp_warming"
)

// Transport wraps http.Client with retry logic and error handling.
type Transport struct {
	httpClient  *http.Client
	apiKey      string
	baseURL     string
	maxRetries  int
	timeout     time.Duration
}

// TransportConfig holds configuration for creating a new Transport.
type TransportConfig struct {
	HTTPClient *http.Client
	APIKey     string
	BaseURL    string
	MaxRetries int
	Timeout    time.Duration
}

// NewTransport creates a new Transport with the given configuration.
// If HTTPClient is nil, a default client will be created.
// If MaxRetries is 0, DefaultMaxRetries will be used.
// If Timeout is 0, DefaultTimeout will be used.
func NewTransport(cfg TransportConfig) *Transport {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	maxRetries := cfg.MaxRetries
	if maxRetries == 0 {
		maxRetries = DefaultMaxRetries
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return &Transport{
		httpClient: httpClient,
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
		maxRetries: maxRetries,
		timeout:    timeout,
	}
}

// Request represents an HTTP request to be made.
type Request struct {
	Method string
	Path   string // URL path (will be appended to baseURL)
	Body   interface{}
	Header http.Header
}

// Response represents an HTTP response.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// Do executes the request with retry logic and returns the response.
// The context controls the overall timeout and cancellation.
func (t *Transport) Do(ctx context.Context, req Request) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		// Check if context is already cancelled
		if ctx.Err() != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return nil, cerebraserrors.NewTimeoutError(ctx)
			}
			return nil, cerebraserrors.NewConnectionError(ctx.Err())
		}

		// Create and execute the request
		resp, err := t.doRequest(ctx, req)
		if err == nil {
			// Success
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable
		if !cerebraserrors.IsRetryableError(err) {
			return nil, err
		}

		// If this was the last attempt, return the error
		if attempt == t.maxRetries {
			return nil, err
		}

		// Calculate backoff duration
		backoff := t.calculateBackoff(attempt, resp)
		
		// Wait before retrying
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return nil, cerebraserrors.NewTimeoutError(ctx)
			}
			return nil, cerebraserrors.NewConnectionError(ctx.Err())
		case <-time.After(backoff):
			// Continue to next retry
		}
	}

	// Should not reach here, but just in case
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, cerebraserrors.NewConnectionError(fmt.Errorf("unexpected retry loop exit"))
}

// doRequest executes a single HTTP request without retry logic.
func (t *Transport) doRequest(ctx context.Context, req Request) (*Response, error) {
	// Build URL
	reqURL, err := url.JoinPath(t.baseURL, req.Path)
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to build URL: %w", err))
	}

	// Create request body
	var bodyReader io.Reader
	if req.Body != nil {
		jsonData, err := json.Marshal(req.Body)
		if err != nil {
			return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to marshal request body: %w", err))
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, reqURL, bodyReader)
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to create request: %w", err))
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+t.apiKey)
	
	// Add custom headers
	if req.Header != nil {
		for key, values := range req.Header {
			for _, value := range values {
				httpReq.Header.Add(key, value)
			}
		}
	}

	// Execute request
	httpResp, err := t.httpClient.Do(httpReq)
	if err != nil {
		// Check for timeout
		if ctx.Err() == context.DeadlineExceeded {
			return nil, cerebraserrors.NewTimeoutError(ctx)
		}
		// Check for URL error (network issues)
		if urlErr, ok := err.(*url.Error); ok {
			return nil, cerebraserrors.NewConnectionError(urlErr)
		}
		return nil, cerebraserrors.NewConnectionError(err)
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to read response body: %w", err))
	}

	// Create response
	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Header:     httpResp.Header.Clone(),
		Body:       body,
	}

	// Check for error status codes
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return resp, cerebraserrors.MapStatusCode(httpResp.StatusCode, string(body))
	}

	return resp, nil
}

// calculateBackoff calculates the backoff duration for a retry attempt.
// Uses exponential backoff with jitter: min(2^attempt * baseBackoff + jitter, maxBackoff)
func (t *Transport) calculateBackoff(attempt int, resp *Response) time.Duration {
	// Check for Retry-After header (for 429 responses)
	if resp != nil && resp.Header != nil {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			// Try to parse as seconds
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				return time.Duration(seconds) * time.Second
			}
			// Try to parse as HTTP date
			if t, err := http.ParseTime(retryAfter); err == nil {
				return time.Until(t)
			}
		}
	}

	// Exponential backoff: 2^attempt * baseBackoff
	backoff := float64(DefaultBaseBackoff)
	for i := 0; i < attempt; i++ {
		backoff *= 2
	}

	// Add jitter (±25%)
	jitter := backoff * DefaultJitterPercent * (rand.Float64()*2 - 1)
	backoff += jitter

	// Cap at max backoff
	if backoff > float64(DefaultMaxBackoff) {
		backoff = float64(DefaultMaxBackoff)
	}

	return time.Duration(backoff)
}

// DoJSON executes a request and unmarshals the JSON response into the target.
func (t *Transport) DoJSON(ctx context.Context, req Request, target interface{}) error {
	resp, err := t.Do(ctx, req)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body, target); err != nil {
		return cerebraserrors.NewConnectionError(fmt.Errorf("failed to unmarshal response: %w", err))
	}

	return nil
}

// WarmTCPConnection performs TCP connection warming to reduce TTFT.
// Returns an error if warming fails, but this should not prevent client construction.
func (t *Transport) WarmTCPConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), TCPWarmingTimeout)
	defer cancel()

	req := Request{
		Method: http.MethodGet,
		Path:   TCPWarmingEndpoint,
	}

	_, err := t.Do(ctx, req)
	return err
}

// GetBaseURL returns the configured base URL.
func (t *Transport) GetBaseURL() string {
	return t.baseURL
}

// GetTimeout returns the configured timeout.
func (t *Transport) GetTimeout() time.Duration {
	return t.timeout
}

// GetHTTPClient returns the underlying HTTP client.
func (t *Transport) GetHTTPClient() *http.Client {
	return t.httpClient
}
