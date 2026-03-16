## Cerebras Cloud SDK for Go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/ldaidone/cerebras-go.svg)](https://pkg.go.dev/github.com/ldaidone/cerebras-cloud-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldaidone/cerebras-go)](https://goreportcard.com/report/github.com/ldaidone/cerebras-cloud-sdk-go)

The Cerebras Cloud Go SDK provides a high-performance, type-safe interface to the Cerebras Cloud API. This library is built for backend reliability, leveraging Go's native concurrency for streaming and strictly avoiding heavy dependencies.

## Installation

```bash
go get github.com/ldaidone/cerebras-cloud-sdk-go
```

## Usage

The SDK mirrors the OpenAI-compatible interface provided by Cerebras, optimized for Go's structural patterns.

### Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ldaidone/cerebras-cloud-sdk-go/pkg/cerebras"
)

func main() {
    // Create client with API key
    client := cerebras.NewClient(
        cerebras.WithAPIKey("sk-your-api-key-here"),
    )

    // Make API calls
    var result map[string]interface{}
    err := client.Get(context.Background(), "/v1/models", &result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Models:", result)
}
```

### Authentication

Pass your API key directly to the client constructor, or let the SDK pull it from the `CEREBRAS_API_KEY` environment variable.

```go
import "github.com/ldaidone/cerebras-cloud-sdk-go/pkg/cerebras"

// Option 1: Explicit API key
client := cerebras.NewClient(
    cerebras.WithAPIKey("your-api-key-here"),
)

// Option 2: Environment variable (set CEREBRAS_API_KEY)
client := cerebras.NewClient()

// Option 3: Explicit key overrides environment
os.Setenv("CEREBRAS_API_KEY", "env-key")
client := cerebras.NewClient(cerebras.WithAPIKey("explicit-key")) // Uses "explicit-key"
```

### Configuration Options

The client supports functional options for flexible configuration:

```go
client := cerebras.NewClient(
    cerebras.WithAPIKey("sk-..."),              // API key (or use CEREBRAS_API_KEY env var)
    cerebras.WithBaseURL("https://api..."),     // Base URL (or use CEREBRAS_BASE_URL env var)
    cerebras.WithTimeout(30*time.Second),       // Request timeout (default: 60s)
    cerebras.WithMaxRetries(5),                 // Max retry attempts (default: 3)
    cerebras.WithTCPWarming(true),              // TCP connection warming (default: true)
    cerebras.WithHTTPClient(customHTTPClient),  // Custom HTTP client (advanced)
)
```

### Local Development

For local development with proxies like Ollama or LM Studio:

```go
client := cerebras.NewClient(
    cerebras.WithBaseURL("http://localhost:11434"),
    cerebras.WithTCPWarming(false),  // Disable warming for local endpoints
)
```

### Error Handling

The SDK provides custom error types for programmatic error handling:

```go
import "github.com/ldaidone/cerebras-cloud-sdk-go/internal/errors"

_, err := client.Get(ctx, "/v1/endpoint", &result)
if err != nil {
    var rateLimitErr *errors.RateLimitError
    if errors.As(err, &rateLimitErr) {
        // Handle rate limiting
        log.Printf("Rate limited, retry after: %v", rateLimitErr.RetryAfter)
    }

    var authErr *errors.AuthenticationError
    if errors.As(err, &authErr) {
        // Handle authentication error
        log.Fatal("Invalid API key")
    }

    if errors.IsRetryableError(err) {
        // Error is eligible for automatic retry
        log.Printf("Retryable error: %v", err)
    }
}
```

**Error Types:**
- `BadRequestError` (400): Invalid request parameters
- `AuthenticationError` (401): Invalid or missing API key
- `PermissionDeniedError` (403): Insufficient permissions
- `NotFoundError` (404): Resource not found
- `UnprocessableEntityError` (422): Invalid request format
- `RateLimitError` (429): Too many requests
- `InternalServerError` (5xx): Server-side errors
- `APIConnectionError`: Network-level errors
- `APITimeoutError`: Request timeout

### Retry Behavior

The SDK automatically retries failed requests with exponential backoff:
- **Default:** 3 retry attempts
- **Backoff:** 2^attempt × 100ms + jitter (±25%)
- **Max backoff:** 2 seconds
- **Retryable errors:** 408, 409, 429, 5xx, network errors, timeouts

### Context Support

All operations support context for timeout and cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := client.Get(ctx, "/v1/endpoint", &result)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
    }
}
```

## Features

### Core Infrastructure
- ✅ **Zero-Dependency Core** - Built on `net/http` (stdlib only)
- ✅ **Context Support** - All operations support `context.Context`
- ✅ **Strict Typing** - Full struct definitions for all API entities
- ✅ **Automatic Retries** - Exponential backoff with jitter (default: 3 retries)
- ✅ **Custom Error Types** - 10 typed errors for programmatic handling
- ✅ **TCP Connection Warming** - Optional TTFT reduction

### API Endpoints
- ✅ **Chat Completions** - Non-streaming and streaming (`/v1/chat/completions`)
- ✅ **Text Completions** - Legacy completions API (`/v1/completions`)
- ✅ **Models API** - List and retrieve models (`/v1/models`)

### Advanced Features
- ✅ **Tool/Function Calling** - Define and call functions from your code
- ✅ **Response Formats** - JSON object and JSON schema constrained output
- ✅ **Service Tiers** - Priority, default, flex, and auto tiers
- ✅ **Reasoning Support** - Reasoning effort and format controls
- ✅ **Raw Response Access** - Access headers and raw response data
- ✅ **Streaming (SSE)** - Channel-based streaming for real-time tokens

### Developer Experience
- ✅ **Functional Options** - Clean, extensible API configuration
- ✅ **Environment Variables** - `CEREBRAS_API_KEY`, `CEREBRAS_BASE_URL`
- ✅ **Model Constants** - Type-safe model identifiers
- ✅ **Helper Functions** - JSON Schema builders, tool helpers
- ✅ **Comprehensive Tests** - 50+ test cases covering all features
- ✅ **GoDoc Documentation** - Full API documentation with examples

## Project Structure

This project follows the [Standard Go Project Layout](PROPOSAL.md):

```
.
├── pkg/cerebras/
│   ├── client.go           # Main client and functional options
│   ├── doc.go              # Package documentation
│   └── example_test.go     # Usage examples
├── internal/
│   ├── transport/
│   │   ├── transport.go    # HTTP transport with retry logic
│   │   └── transport_test.go
│   └── errors/
│       ├── errors.go       # Custom error types
│       └── errors_test.go
├── openspec/
│   ├── changes/            # OpenSpec change artifacts
│   └── specs/              # API specifications
├── go.mod
├── README.md
└── PROPOSAL.md
```

## Contributing

We welcome contributions! Please see our [Spec-Driven Development](https://www.google.com/search?q=./PROPOSAL.md) document for architectural decisions regarding this port.

## License

Apache License 2.0
