## Why

The Cerebras Cloud SDK needs a Go client that provides a high-performance, type-safe interface to the Cerebras Cloud API. The official Python SDK exists, but Go offers superior performance for backend services, native concurrency patterns for streaming, and smaller binary sizes. This change establishes the foundation: client initialization with configuration, authentication, built-in retry logic, and the HTTP transport layer.

## What Changes

- New `Client` struct with functional options pattern (`WithAPIKey`, `WithBaseURL`, `WithTimeout`, `WithMaxRetries`, etc.)
- Support for API key authentication (direct parameter or `CEREBRAS_API_KEY` env var)
- Configurable base URL for local AI parity (Ollama, LM Studio, private endpoints)
- Built-in retry strategy with exponential backoff (Go-idiomatic defaults)
- Timeout support via `context.Context` on all operations
- Optional TCP connection warming on client construction
- Custom error types: `APIError`, `RateLimitError`, `AuthenticationError`, `BadRequestError`, `APIConnectionError`, `APITimeoutError`
- HTTP transport wrapper with automatic retry logic for transient failures

## Capabilities

### New Capabilities
- `client-initialization`: Core client setup, configuration, authentication, and functional options pattern
- `error-handling`: Custom error types, error hierarchy, and retry behavior specification
- `http-transport`: HTTP client wrapper with retries, timeouts, connection pooling, and TCP warming

### Modified Capabilities
- *(none - this is the initial implementation)*

## Impact

- Creates the foundation for all subsequent API calls (chat completions, streaming, etc.)
- Establishes the functional options pattern used throughout the SDK
- Defines error handling conventions and retry semantics
- Introduces the `internal/transport` package for HTTP abstraction
- No breaking changes (greenfield project)
- Dependencies: Standard library only (`net/http`, `context`, `time`, `os`)
