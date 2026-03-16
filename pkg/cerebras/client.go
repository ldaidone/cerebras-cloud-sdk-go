// Package cerebras provides a high-performance, type-safe Go client for the Cerebras Cloud API.
//
// This SDK mirrors the functionality of the official Python SDK while leveraging Go's strengths:
//   - Performance and concurrency
//   - Minimal dependencies (standard library only)
//   - Native streaming support via channels
//   - Context-aware timeouts and cancellation
//
// Basic usage:
//
//	client := cerebras.NewClient(
//	    cerebras.WithAPIKey("your-api-key"),
//	)
//
// The client supports functional options for configuration:
//
//	client := cerebras.NewClient(
//	    cerebras.WithAPIKey("your-api-key"),
//	    cerebras.WithBaseURL("https://api.cerebras.ai"),
//	    cerebras.WithTimeout(30 * time.Second),
//	    cerebras.WithMaxRetries(5),
//	)
//
// API keys can also be set via the CEREBRAS_API_KEY environment variable.
package cerebras

import (
	"context"
	"net/http"
	"os"
	"time"

	"github/ldaidone/cerebras-cloud-sdk-go/internal/transport"
)

// Default configuration values.
const (
	DefaultBaseURL    = "https://api.cerebras.ai"
	DefaultTimeout    = 60 * time.Second
	DefaultMaxRetries = 3
)

// Client is the main entry point for the Cerebras Cloud SDK.
// It provides methods for interacting with the Cerebras Cloud API.
//
// The client is safe for concurrent use by multiple goroutines.
type Client struct {
	apiKey           string
	baseURL          string
	timeout          time.Duration
	maxRetries       int
	tcpWarming       bool
	httpClient       *http.Client
	transport        *transport.Transport
	ChatCompletions  *ChatCompletionsService
	Models           *ModelsService
	TextCompletions  *TextCompletionsService
}

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithAPIKey sets the API key for authentication.
// This takes precedence over the CEREBRAS_API_KEY environment variable.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithBaseURL sets the base URL for API requests.
// This is useful for local development with proxies like Ollama or LM Studio.
// Default: https://api.cerebras.ai
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout sets the default timeout for API requests.
// Individual requests can override this using context.WithTimeout.
// Default: 60 seconds
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

// WithMaxRetries sets the maximum number of retry attempts for failed requests.
// Set to 0 to disable retries.
// Default: 3 retries
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		c.maxRetries = n
	}
}

// WithTCPWarming enables or disables TCP connection warming on client construction.
// When enabled, the client sends a warming request to reduce time-to-first-token (TTFT).
// Warming failures are logged but do not prevent client construction.
// Default: true (enabled)
func WithTCPWarming(enabled bool) Option {
	return func(c *Client) {
		c.tcpWarming = enabled
	}
}

// WithHTTPClient sets a custom HTTP client.
// This is an escape hatch for advanced use cases.
// If not set, a default http.Client will be created.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient creates a new Cerebras Cloud API client.
//
// Configuration options:
//   - WithAPIKey: Sets the API key (or use CEREBRAS_API_KEY env var)
//   - WithBaseURL: Sets the base URL (or use CEREBRAS_BASE_URL env var)
//   - WithTimeout: Sets the request timeout (default: 60s)
//   - WithMaxRetries: Sets max retry attempts (default: 3)
//   - WithTCPWarming: Enables TCP warming (default: true)
//   - WithHTTPClient: Uses a custom HTTP client
//
// Example:
//
//	client := cerebras.NewClient(
//	    cerebras.WithAPIKey("sk-..."),
//	    cerebras.WithTimeout(30 * time.Second),
//	)
func NewClient(opts ...Option) *Client {
	// Initialize with defaults
	c := &Client{
		baseURL:    DefaultBaseURL,
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
		tcpWarming: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Fall back to environment variables if not set
	if c.apiKey == "" {
		c.apiKey = os.Getenv("CEREBRAS_API_KEY")
	}
	if c.baseURL == DefaultBaseURL {
		if envURL := os.Getenv("CEREBRAS_BASE_URL"); envURL != "" {
			c.baseURL = envURL
		}
	}

	// Create transport layer
	c.transport = transport.NewTransport(transport.TransportConfig{
		HTTPClient: c.httpClient,
		APIKey:     c.apiKey,
		BaseURL:    c.baseURL,
		MaxRetries: c.maxRetries,
		Timeout:    c.timeout,
	})

	// Initialize services
	c.ChatCompletions = &ChatCompletionsService{
		client:    c,
		transport: c.transport,
	}
	c.Models = &ModelsService{
		client:    c,
		transport: c.transport,
	}
	c.TextCompletions = &TextCompletionsService{
		client:    c,
		transport: c.transport,
	}

	// Perform TCP warming if enabled
	if c.tcpWarming {
		if err := c.transport.WarmTCPConnection(); err != nil {
			// Log warming failure but don't fail construction
			// Using standard Go logging - in production, use slog
			// log.Printf("Warning: TCP warming failed: %v", err)
		}
	}

	return c
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Timeout returns the configured timeout.
func (c *Client) Timeout() time.Duration {
	return c.timeout
}

// Do executes a generic HTTP request using the transport layer.
// This is a low-level method for internal use and advanced scenarios.
//
// For most use cases, use the higher-level methods like CreateChatCompletion.
func (c *Client) Do(ctx context.Context, method string, path string, body interface{}, result interface{}) error {
	req := transport.Request{
		Method: method,
		Path:   path,
		Body:   body,
	}

	return c.transport.DoJSON(ctx, req, result)
}

// Get executes a GET request.
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Do(ctx, http.MethodGet, path, nil, result)
}

// Post executes a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Do(ctx, http.MethodPost, path, body, result)
}
