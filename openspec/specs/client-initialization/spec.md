# Capability: Client Initialization

The client initialization capability provides the foundation for creating and configuring the Cerebras Cloud SDK client.

## Requirements

### Requirement: Client initialization with functional options
The system SHALL provide a `Client` struct that can be configured using functional options pattern.

#### Scenario: Create client with default configuration
- **WHEN** user calls `NewClient()` with no options
- **THEN** system creates a client with default values (API key from env var, default base URL, default timeout, default retries)

#### Scenario: Create client with API key
- **WHEN** user calls `NewClient(WithAPIKey("sk-..."))`
- **THEN** system creates a client configured with the provided API key

#### Scenario: Create client with custom base URL
- **WHEN** user calls `NewClient(WithBaseURL("https://api.example.com"))`
- **THEN** system creates a client configured to use the custom base URL

#### Scenario: Create client with custom timeout
- **WHEN** user calls `NewClient(WithTimeout(30 * time.Second))`
- **THEN** system creates a client with the specified timeout for all requests

#### Scenario: Create client with multiple options
- **WHEN** user calls `NewClient(WithAPIKey("sk-..."), WithBaseURL("..."), WithTimeout(30*time.Second))`
- **THEN** system creates a client with all specified options applied

#### Scenario: API key from environment variable
- **WHEN** user calls `NewClient()` without `WithAPIKey` and `CEREBRAS_API_KEY` env var is set
- **THEN** system uses the API key from the environment variable

#### Scenario: Explicit API key overrides environment
- **WHEN** user calls `NewClient(WithAPIKey("sk-..."))` and `CEREBRAS_API_KEY` env var is also set
- **THEN** system uses the explicitly provided API key, not the environment variable

### Requirement: TCP connection warming
The system SHALL optionally perform TCP connection warming on client construction to reduce time-to-first-token (TTFT).

#### Scenario: TCP warming enabled by default
- **WHEN** user creates a client without specifying `WithTCPWarming`
- **THEN** system performs TCP warming by sending a request to `/v1/tcp_warming`

#### Scenario: TCP warming disabled
- **WHEN** user calls `NewClient(WithTCPWarming(false))`
- **THEN** system skips TCP warming on construction

#### Scenario: TCP warming failure is non-fatal
- **WHEN** TCP warming request fails
- **THEN** client construction succeeds and warming failure is logged (not returned as error)

### Requirement: Context support
The system SHALL accept `context.Context` as the first parameter in all public methods.

#### Scenario: Request respects context cancellation
- **WHEN** user passes a cancelled context to a client method
- **THEN** system returns an error immediately without making the request

#### Scenario: Request respects context timeout
- **WHEN** user passes a context with timeout and request exceeds it
- **THEN** system cancels the request and returns a timeout error

### Requirement: Client configuration accessors
The system SHALL provide methods to access current client configuration.

#### Scenario: Get current base URL
- **WHEN** user calls `client.BaseURL()`
- **THEN** system returns the configured base URL

#### Scenario: Get current timeout
- **WHEN** user calls `client.Timeout()`
- **THEN** system returns the configured timeout duration
