## 1. Project Structure

- [x] 1.1 Create directory structure: `pkg/cerebras/`, `internal/transport/`, `internal/errors/`
- [x] 1.2 Initialize Go module (if not already done)
- [x] 1.3 Create package documentation comments

## 2. Error Types Implementation

- [x] 2.1 Define base `APIError` interface in `internal/errors/errors.go`
- [x] 2.2 Implement `APIStatusError` with status code field
- [x] 2.3 Implement specific status errors: `BadRequestError`, `AuthenticationError`, `PermissionDeniedError`, `NotFoundError`, `UnprocessableEntityError`, `RateLimitError`, `InternalServerError`
- [x] 2.4 Implement `APIConnectionError` and `APITimeoutError`
- [x] 2.5 Add `IsRetryableError(err error)` helper function
- [x] 2.6 Implement `Error()` methods for all error types with descriptive messages
- [x] 2.7 Add `Unwrap()` methods for errors.Is/As support

## 3. HTTP Transport Layer

- [x] 3.1 Create `internal/transport/transport.go` with `Transport` struct
- [x] 3.2 Implement `NewTransport()` with default HTTP client configuration
- [x] 3.3 Implement `Do()` method with retry logic
- [x] 3.4 Implement exponential backoff calculation with jitter
- [x] 3.5 Add retry condition checking (which errors trigger retry)
- [x] 3.6 Implement Retry-After header parsing and respect
- [x] 3.7 Add max backoff ceiling (2 seconds default)
- [x] 3.8 Implement request timeout handling from context
- [x] 3.9 Implement JSON request body marshaling
- [x] 3.10 Implement JSON response unmarshaling
- [x] 3.11 Implement response status code to error type mapping

## 4. Client Implementation

- [x] 4.1 Create `pkg/cerebras/client.go` with `Client` struct
- [x] 4.2 Define `Option` type as `func(*Client)`
- [x] 4.3 Implement `WithAPIKey()` option
- [x] 4.4 Implement `WithBaseURL()` option
- [x] 4.5 Implement `WithTimeout()` option
- [x] 4.6 Implement `WithMaxRetries()` option
- [x] 4.7 Implement `WithTCPWarming()` option
- [x] 4.8 Implement `WithHTTPClient()` option for custom client injection
- [x] 4.9 Implement `NewClient()` variadic constructor
- [x] 4.10 Add environment variable fallback for API key (`CEREBRAS_API_KEY`)
- [x] 4.11 Add environment variable fallback for base URL (`CEREBRAS_BASE_URL`)
- [x] 4.12 Set default values (timeout: 60s, maxRetries: 3, tcpWarming: true)

## 5. TCP Connection Warming

- [x] 5.1 Implement `warmTCPConnection()` method in client
- [x] 5.2 Send GET request to `<baseURL>/v1/tcp_warming` on construction
- [x] 5.3 Use separate short-lived context (5 second timeout) for warming
- [x] 5.4 Log warming failures without failing construction (use `slog` or configurable logger)

## 6. Configuration Accessors

- [x] 6.1 Implement `BaseURL()` method to return configured base URL
- [x] 6.2 Implement `Timeout()` method to return configured timeout

## 7. Integration

- [x] 7.1 Wire up client to use transport layer for all HTTP operations
- [x] 7.2 Ensure context is passed through all layers (client → transport → HTTP)
- [x] 7.3 Add `Authorization: Bearer <api-key>` header to all requests
- [x] 7.4 Implement base URL + endpoint URL construction

## 8. Testing

- [x] 8.1 Create unit tests for error types and `IsRetryableError()`
- [x] 8.2 Create unit tests for exponential backoff calculation
- [x] 8.3 Create unit tests for retry condition checking
- [x] 8.4 Create unit tests for client options
- [x] 8.5 Create integration test with mocked HTTP client
- [x] 8.6 Create test for TCP warming behavior
- [x] 8.7 Create test for environment variable fallback
- [x] 8.8 Create test for error type mapping from HTTP status codes

## 9. Documentation

- [x] 9.1 Add Go doc comments to all public types and functions
- [x] 9.2 Create example usage in `pkg/cerebras/doc.go`
- [x] 9.3 Add README section for client initialization
- [x] 9.4 Create example file showing all configuration options
