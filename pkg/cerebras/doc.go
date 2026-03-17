// Package cerebras provides a high-performance, type-safe Go client for the Cerebras Cloud API.
//
// # Getting Started
//
// Install the SDK:
//
//	go get github.com/ldaidone/cerebras-cloud-sdk-go
//
// Create a client:
//
//	import "github.com/ldaidone/cerebras-cloud-sdk-go/pkg/cerebras"
//
//	client := cerebras.NewClient(
//	    cerebras.WithAPIKey("your-api-key"),
//	)
//
// # Authentication
//
// The API key can be provided in several ways:
//
// 1. Explicitly via WithAPIKey option:
//
//	client := cerebras.NewClient(cerebras.WithAPIKey("sk-..."))
//
// 2. Via environment variable:
//
//	export CEREBRAS_API_KEY=sk-...
//	client := cerebras.NewClient()
//
// 3. Explicit key takes precedence over environment variable.
//
// # Configuration
//
// The client supports various configuration options:
//
//	client := cerebras.NewClient(
//	    cerebras.WithAPIKey("sk-..."),
//	    cerebras.WithBaseURL("https://api.cerebras.ai"),  // Or custom endpoint
//	    cerebras.WithTimeout(30 * time.Second),           // Request timeout
//	    cerebras.WithMaxRetries(5),                       // Retry attempts
//	    cerebras.WithTCPWarming(true),                    // TCP connection warming
//	)
//
// # Environment Variables
//
// - CEREBRAS_API_KEY: Default API key (used if WithAPIKey is not provided)
// - CEREBRAS_BASE_URL: Default base URL (used if WithBaseURL is not provided)
//
// # Local Development
//
// For local development with proxies like Ollama or LM Studio:
//
//	client := cerebras.NewClient(
//	    cerebras.WithBaseURL("http://localhost:11434"),
//	    cerebras.WithTCPWarming(false),  // Disable warming for local endpoints
//	)
//
// # Error Handling
//
// The SDK provides custom error types for programmatic error handling:
//
//	import "github.com/ldaidone/cerebras-cloud-sdk-go/internal/errors"
//
//	_, err := client.Get(ctx, "/v1/some-endpoint", &result)
//	if err != nil {
//	    var rateLimitErr *errors.RateLimitError
//	    if errors.As(err, &rateLimitErr) {
//	        // Handle rate limiting
//	        log.Printf("Rate limited, retry after: %v", rateLimitErr.RetryAfter)
//	    }
//
//	    var authErr *errors.AuthenticationError
//	    if errors.As(err, &authErr) {
//	        // Handle authentication error
//	        log.Fatal("Invalid API key")
//	    }
//
//	    if errors.IsRetryableError(err) {
//	        // Error is eligible for automatic retry
//	        log.Printf("Retryable error: %v", err)
//	    }
//	}
//
// # Error Types
//
// - BadRequestError (400): Invalid request parameters
// - AuthenticationError (401): Invalid or missing API key
// - PermissionDeniedError (403): Insufficient permissions
// - NotFoundError (404): Resource not found
// - UnprocessableEntityError (422): Invalid request format
// - RateLimitError (429): Too many requests
// - InternalServerError (5xx): Server-side errors
// - APIConnectionError: Network-level errors
// - APITimeoutError: Request timeout
//
// # Retry Behavior
//
// The SDK automatically retries failed requests with exponential backoff:
//
//   - Default: 3 retry attempts
//   - Backoff: 2^attempt * 100ms + jitter (±25%)
//   - Max backoff: 2 seconds
//   - Retryable errors: 408, 409, 429, 5xx, network errors, timeouts
//
// # Context Support
//
// All operations support context for timeout and cancellation:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	err := client.Get(ctx, "/v1/some-endpoint", &result)
//	if err != nil {
//	    if errors.Is(err, context.DeadlineExceeded) {
//	        // Handle timeout
//	    }
//	}
//
// # TCP Connection Warming
//
// By default, the client performs TCP connection warming on construction to reduce
// time-to-first-token (TTFT). This sends a preliminary request to establish the
// connection before your actual API calls.
//
// Disable for local development or air-gapped environments:
//
//	client := cerebras.NewClient(cerebras.WithTCPWarming(false))
//
// # Custom HTTP Client
//
// For advanced use cases, you can provide a custom HTTP client:
//
//	customClient := &http.Client{
//	    Timeout: 30 * time.Second,
//	    Transport: &http.Transport{
//	        MaxIdleConns:        100,
//	        MaxIdleConnsPerHost: 10,
//	    },
//	}
//
//	client := cerebras.NewClient(
//	    cerebras.WithAPIKey("sk-..."),
//	    cerebras.WithHTTPClient(customClient),
//	)
//
// # Thread Safety
//
// The Client is safe for concurrent use by multiple goroutines.
package cerebras
