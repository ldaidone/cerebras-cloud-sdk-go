# Capability: Error Handling

The error handling capability provides a comprehensive error hierarchy for programmatic error handling in the Cerebras Cloud SDK.

## Requirements

### Requirement: Error hierarchy
The system SHALL provide a hierarchy of custom error types for programmatic error handling.

#### Scenario: Base APIError type
- **WHEN** an API-related error occurs
- **THEN** system returns an error that implements the `APIError` interface

#### Scenario: HTTP status errors
- **WHEN** API returns a 4xx or 5xx status code
- **THEN** system returns an `APIStatusError` with the HTTP status code and response body

#### Scenario: BadRequestError for 400
- **WHEN** API returns HTTP 400
- **THEN** system returns a `BadRequestError`

#### Scenario: AuthenticationError for 401
- **WHEN** API returns HTTP 401
- **THEN** system returns an `AuthenticationError`

#### Scenario: PermissionDeniedError for 403
- **WHEN** API returns HTTP 403
- **THEN** system returns a `PermissionDeniedError`

#### Scenario: NotFoundError for 404
- **WHEN** API returns HTTP 404
- **THEN** system returns a `NotFoundError`

#### Scenario: UnprocessableEntityError for 422
- **WHEN** API returns HTTP 422
- **THEN** system returns an `UnprocessableEntityError`

#### Scenario: RateLimitError for 429
- **WHEN** API returns HTTP 429
- **THEN** system returns a `RateLimitError`

#### Scenario: InternalServerError for 5xx
- **WHEN** API returns HTTP 500, 502, 503, or 504
- **THEN** system returns an `InternalServerError`

#### Scenario: Connection errors
- **WHEN** a network error occurs (DNS failure, connection refused, TLS error)
- **THEN** system returns an `APIConnectionError`

#### Scenario: Timeout errors
- **WHEN** a request times out (context.DeadlineExceeded)
- **THEN** system returns an `APITimeoutError`

### Requirement: Error type inspection
The system SHALL provide ways to inspect and categorize errors programmatically.

#### Scenario: Type assertion to specific error
- **WHEN** user performs type assertion `err.(*RateLimitError)`
- **THEN** system allows access to rate limit-specific fields (retry-after, etc.)

#### Scenario: errors.Is support
- **WHEN** user calls `errors.Is(err, &RateLimitError{})`
- **THEN** system returns true if error is a rate limit error

#### Scenario: errors.As support
- **WHEN** user calls `errors.As(err, &targetErr)` where targetErr is `*APIStatusError`
- **THEN** system populates targetErr if the error is an API status error

#### Scenario: Error message format
- **WHEN** an error is converted to string (Error() method)
- **THEN** system returns a human-readable message including status code (if applicable) and description

### Requirement: Retry eligibility
The system SHALL identify which errors are eligible for automatic retry.

#### Scenario: Rate limit errors are retryable
- **WHEN** a `RateLimitError` occurs
- **THEN** system marks the error as retryable

#### Scenario: Server errors are retryable
- **WHEN** an `InternalServerError` occurs (500, 502, 503, 504)
- **THEN** system marks the error as retryable

#### Scenario: Connection errors are retryable
- **WHEN** an `APIConnectionError` occurs
- **THEN** system marks the error as retryable

#### Scenario: Client errors are not retryable
- **WHEN** a `BadRequestError`, `AuthenticationError`, or `PermissionDeniedError` occurs
- **THEN** system marks the error as NOT retryable

#### Scenario: Helper function for retry check
- **WHEN** user calls `IsRetryableError(err)`
- **THEN** system returns true if the error type is eligible for retry
