## Context

This is a greenfield Go SDK for the Cerebras Cloud API, mirroring the functionality of the official Python SDK while leveraging Go's strengths: performance, concurrency, and minimal dependencies. The Python SDK uses `httpx`, Pydantic models, and generator-based streaming. This Go port will use `net/http`, structs with JSON tags, and channel-based streaming.

**Constraints:**
- Zero external dependencies in core (standard library only)
- Must support both streaming and non-streaming responses
- Must handle retries transparently with exponential backoff
- Must support context-based cancellation and timeouts
- API must be idiomatic Go (functional options, error types, channels)

**Reference:** Cerebras Cloud SDK for Python - https://github.com/Cerebras/cerebras-cloud-sdk-python

## Goals / Non-Goals

**Goals:**
- Provide a clean, type-safe client initialization API with functional options
- Implement automatic retry logic with configurable exponential backoff
- Support TCP connection warming to reduce time-to-first-token (TTFT)
- Define a clear error hierarchy for programmatic error handling
- Establish the foundation for chat completions (streaming and non-streaming)
- Keep binary size minimal and supply chain secure (no external deps)

**Non-Goals:**
- Actual API endpoint implementations (chat completions, etc.) - covered in subsequent changes
- Async/await patterns (Go uses goroutines, not async/await)
- Custom HTTP client injection (users can override via `WithHTTPClient` if needed, but default is `net/http`)
- Middleware or interceptor patterns (keep it simple for v1)

## Decisions

### 1. Client Structure: Single Client with Options

**Decision:** Single `Client` struct with functional options, not separate sync/async clients.

**Rationale:** Go doesn't need separate async clients - concurrency is handled via goroutines and channels. The Python SDK has `Cerebras` and `AsyncCerebras` due to Python's sync/async divide. In Go, one client handles both patterns naturally.

**Alternatives Considered:**
- Separate `SyncClient` and `AsyncClient`: Adds complexity without benefit in Go
- Builder pattern: Functional options are more idiomatic in Go and easier to extend

```go
type Client struct {
    apiKey      string
    baseURL     string
    timeout     time.Duration
    maxRetries  int
    httpClient  *http.Client
    // ...
}

func NewClient(opts ...Option) *Client
```

### 2. Functional Options Pattern

**Decision:** Use functional options for configuration.

**Rationale:** Prevents breaking changes when adding new options, provides a clean API, and is idiomatic Go.

```go
type Option func(*Client)

func WithAPIKey(key string) Option
func WithBaseURL(url string) Option
func WithTimeout(d time.Duration) Option
func WithMaxRetries(n int) Option
func WithTCPWarming(enabled bool) Option
```

**Alternatives Considered:**
- Config struct: Harder to extend, requires updating all call sites
- Builder pattern: More verbose, less flexible for optional params

### 3. Error Hierarchy

**Decision:** Custom error types matching the Python SDK's error hierarchy.

**Rationale:** Users need to programmatically handle different error scenarios (rate limits, auth failures, network issues). Mirroring the Python SDK makes migration easier.

```
APIError (base)
├── APIStatusError (HTTP status errors)
│   ├── BadRequestError (400)
│   ├── AuthenticationError (401)
│   ├── PermissionDeniedError (403)
│   ├── NotFoundError (404)
│   ├── UnprocessableEntityError (422)
│   ├── RateLimitError (429)
│   └── InternalServerError (>=500)
└── APIConnectionError (network errors)
    └── APITimeoutError (context.DeadlineExceeded)
```

**Alternatives Considered:**
- Wrap all errors in a single type with status code field: Loses type safety
- Use only standard errors: Can't distinguish error types without string parsing

### 4. Retry Strategy: Exponential Backoff with Jitter

**Decision:** Built-in retry with exponential backoff, defaulting to 3 retries (Go performance expectations).

**Rationale:** The Python SDK uses 2 retries. Go services typically handle higher throughput, so 3 retries provides better resilience without excessive latency. Backoff formula: `min(time.Duration(2^attempt) * 100ms + jitter, maxBackoff)`

**Retry Conditions:**
- Connection errors
- 408 Request Timeout
- 409 Conflict
- 429 Rate Limit
- 5xx Server Errors

**Alternatives Considered:**
- No retry (user handles): Puts burden on users, poor UX
- Fixed retry delay: Less resilient to burst failures
- 2 retries (Python default): Conservative for Go workloads

### 5. TCP Connection Warming

**Decision:** Optional TCP warming on client construction, enabled by default.

**Rationale:** The Python SDK sends a warming request to `/v1/tcp_warming` on construction to reduce TTFT. This is valuable for production but should be disableable for local development or air-gapped environments.

```go
func WithTCPWarming(enabled bool) Option  // default: true
```

**Alternatives Considered:**
- Always warm: Inflexible for local dev
- Never warm: Leaves performance on the table
- Separate method call: Less discoverable

### 6. HTTP Transport Layer

**Decision:** Internal `transport` package wraps `http.Client` with retry logic.

**Rationale:** Separates concerns - client handles configuration, transport handles HTTP mechanics. Makes testing easier and allows future enhancements (circuit breakers, metrics) without changing client API.

```
┌─────────────────────────────────────────────────────────────┐
│                      Client                                  │
│  - Configuration (apiKey, baseURL, timeout, retries)        │
│  - Public API surface                                        │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ delegates to
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Transport                                  │
│  - Retry logic with exponential backoff                     │
│  - Request/response transformation                          │
│  - Error mapping to custom error types                      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ wraps
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   http.Client                                │
│  - Actual HTTP communication                                │
│  - Connection pooling                                        │
│  - TLS handling                                              │
└─────────────────────────────────────────────────────────────┘
```

**Alternatives Considered:**
- Inline HTTP logic in Client: Harder to test, violates separation of concerns
- Interface-based transport: Over-engineering for v1, can add later if needed

### 7. Context Propagation

**Decision:** All public methods accept `context.Context` as first parameter.

**Rationale:** Idiomatic Go for timeout and cancellation support. Users can control request lifetime, and it integrates with Go's cancellation propagation.

```go
func (c *Client) DoSomething(ctx context.Context, ...) (..., error)
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Exponential backoff may increase latency for transient failures | Cap max backoff at reasonable value (e.g., 2s), add jitter to prevent thundering herd |
| TCP warming adds startup overhead | Make it optional, document trade-off clearly |
| Custom error types increase API surface | Provide helper functions (`IsRateLimit`, `IsAuthError`) for common checks |
| No external HTTP client means missing features (circuit breakers, metrics) | Provide `WithHTTPClient` escape hatch, can add features in v2 |
| Channel-based streaming (Option C) requires careful channel management | Document patterns clearly, provide examples, use buffered channels appropriately |

## Open Questions

1. **Default timeout value**: Python SDK uses 60s. Is this appropriate for Go, or should we use a shorter default (e.g., 30s) with easy override?

2. **Max backoff ceiling**: What's a reasonable max backoff? 2s? 5s? Should it be configurable?

3. **TCP warming endpoint**: The Python SDK uses `/v1/tcp_warming`. Is this the correct endpoint for Cerebras Cloud, or is this internal to their SDK?

4. **Error wrapping**: Should we wrap underlying `net/http` errors or translate them completely? Go 1.13+ error wrapping (`%w`) vs. full translation.
