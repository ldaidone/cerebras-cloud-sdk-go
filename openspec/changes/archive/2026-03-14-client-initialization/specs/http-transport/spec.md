## ADDED Requirements

### Requirement: HTTP client wrapper
The system SHALL provide an internal HTTP transport layer that wraps `net/http.Client`.

#### Scenario: Default HTTP client creation
- **WHEN** client is created without `WithHTTPClient` option
- **THEN** system creates a default `http.Client` with appropriate settings

#### Scenario: Custom HTTP client injection
- **WHEN** user calls `NewClient(WithHTTPClient(customClient))`
- **THEN** system uses the provided HTTP client instead of creating a default

#### Scenario: Default client configuration
- **WHEN** default HTTP client is created
- **THEN** system configures reasonable defaults for connection pooling and TLS

### Requirement: Automatic retry with exponential backoff
The system SHALL automatically retry failed requests using exponential backoff with jitter.

#### Scenario: Default retry count
- **WHEN** client is created without `WithMaxRetries` option
- **THEN** system uses 3 as the default maximum retry count

#### Scenario: Custom retry count
- **WHEN** user calls `NewClient(WithMaxRetries(5))`
- **THEN** system retries failed requests up to 5 times

#### Scenario: Exponential backoff calculation
- **WHEN** a retryable error occurs on attempt N
- **THEN** system waits `min(2^N * 100ms + jitter, maxBackoff)` before retrying

#### Scenario: Jitter prevents thundering herd
- **WHEN** multiple retries occur
- **THEN** system adds random jitter (±25% of backoff duration) to each retry delay

#### Scenario: Max backoff ceiling
- **WHEN** backoff calculation exceeds max backoff (2 seconds)
- **THEN** system caps the wait time at the maximum

#### Scenario: Retry on connection error
- **WHEN** a request fails with `APIConnectionError`
- **THEN** system retries the request (up to max retries)

#### Scenario: Retry on rate limit
- **WHEN** a request fails with `RateLimitError` (429)
- **THEN** system retries the request, respecting Retry-After header if present

#### Scenario: Retry on server error
- **WHEN** a request fails with `InternalServerError` (5xx)
- **THEN** system retries the request (up to max retries)

#### Scenario: No retry on client error
- **WHEN** a request fails with `BadRequestError` or `AuthenticationError` (4xx, non-retryable)
- **THEN** system returns the error immediately without retrying

#### Scenario: Exhausted retries
- **WHEN** all retry attempts are exhausted
- **THEN** system returns the last error to the caller

#### Scenario: Retry-After header respected
- **WHEN** a 429 response includes a `Retry-After` header
- **THEN** system waits for the specified duration before retrying (overrides backoff calculation)

### Requirement: Request timeout handling
The system SHALL enforce request timeouts at the HTTP transport layer.

#### Scenario: Timeout from context
- **WHEN** request context has a deadline
- **THEN** system respects the context deadline for the HTTP request

#### Scenario: Default timeout
- **WHEN** no context deadline is set and client has default timeout
- **THEN** system applies the client's configured timeout (default: 60 seconds)

#### Scenario: Timeout error mapping
- **WHEN** a request times out
- **THEN** system returns an `APITimeoutError` (not a generic context error)

### Requirement: Request/Response handling
The system SHALL handle HTTP request creation and response parsing.

#### Scenario: JSON request body
- **WHEN** request includes a body struct
- **THEN** system marshals the body to JSON and sets `Content-Type: application/json`

#### Scenario: JSON response parsing
- **WHEN** response has `Content-Type: application/json`
- **THEN** system unmarshals the response body into the target struct

#### Scenario: API key header injection
- **WHEN** making any request
- **THEN** system adds `Authorization: Bearer <api-key>` header

#### Scenario: Base URL prepending
- **WHEN** making a request to endpoint `/v1/chat/completions`
- **THEN** system constructs full URL as `<baseURL>/v1/chat/completions`

#### Scenario: Response status code checking
- **WHEN** response is received
- **THEN** system checks status code and maps to appropriate error type if not 2xx

### Requirement: TCP connection warming
The system SHALL implement TCP connection warming to reduce TTFT.

#### Scenario: Warming request endpoint
- **WHEN** TCP warming is enabled
- **THEN** system sends a GET request to `<baseURL>/v1/tcp_warming` on client construction

#### Scenario: Warming uses separate context
- **WHEN** performing TCP warming
- **THEN** system uses a separate short-lived context (5 second timeout) independent of user contexts

#### Scenario: Warming failure handling
- **WHEN** TCP warming request fails
- **THEN** system logs the error (if logging is configured) but does not fail client construction
