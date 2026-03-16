## Context

The Go SDK has completed its foundation layer, type definitions, chat completions, streaming, tool calling, response formats, and text completions. Advanced features including service tiers, reasoning model support, and raw response access provide additional capabilities for power users and debugging scenarios. These features match the Python SDK's advanced functionality.

Key constraints:
- API compatibility with Python SDK advanced features
- Non-breaking additions to existing types
- Clean integration with existing client architecture
- Minimal complexity for common use cases

## Goals / Non-Goals

**Goals:**
- Add service_tier parameter (auto, flex, priority)
- Add reasoning_effort and reasoning_format parameters
- Add clear_thinking and disable_reasoning parameters
- Implement raw response access pattern
- Add HTTP verb helpers for custom API calls
- Add WithOptions method for client configuration
- Comprehensive examples and documentation

**Non-Goals:**
- Custom retry logic configuration (existing retry is sufficient)
- Custom transport layer (users can provide custom HTTP client)
- Advanced authentication methods (API key is standard)

## Decisions

### 1. Service Tier Type
**Decision:** Define ServiceTier as string enum with constants

**Rationale:**
- Type safety without external dependencies
- Clear IDE autocomplete support
- Easy JSON marshaling via string conversion

**Structure:**
```go
type ServiceTier string
const (
    ServiceTierAuto     ServiceTier = "auto"
    ServiceTierFlex     ServiceTier = "flex"
    ServiceTierPriority ServiceTier = "priority"
)
```

### 2. Reasoning Parameters
**Decision:** Add reasoning_effort, reasoning_format, clear_thinking, disable_reasoning as optional fields

**Rationale:**
- Matches Cerebras API specification
- Optional fields don't break existing code
- Pointer types distinguish unset from zero values

**Structure:**
```go
type ChatCompletionRequest struct {
    // ... existing fields
    ReasoningEffort   *string `json:"reasoning_effort,omitempty"`
    ReasoningFormat   *string `json:"reasoning_format,omitempty"`
    ClearThinking     *bool   `json:"clear_thinking,omitempty"`
    DisableReasoning  *bool   `json:"disable_reasoning,omitempty"`
}
```

### 3. Reasoning Effort Values
**Decision:** Define constants for reasoning effort levels

**Rationale:**
- Prevents typos in effort values
- Documents supported values
- IDE autocomplete

**Constants:**
```go
const (
    ReasoningEffortLow    string = "low"
    ReasoningEffortMedium string = "medium"
    ReasoningEffortHigh   string = "high"
)
```

### 4. Raw Response Access Pattern
**Decision:** Implement wrapper type with response metadata access

**Rationale:**
- Python SDK uses `.with_raw_response` pattern
- Go can use method-based approach
- Provides access to status code, headers, and raw body

**Structure:**
```go
type APIResponse[T any] struct {
    Data       T
    StatusCode int
    Headers    http.Header
    RawBody    []byte
}

func (c *Client) ChatCompletions.CreateWithResponse(ctx context.Context, req *ChatCompletionRequest) (*APIResponse[ChatCompletion], error)
```

### 5. HTTP Verb Helpers
**Decision:** Add generic HTTP method helpers to client

**Rationale:**
- Enables custom API calls without implementing every endpoint
- Reuses existing transport layer
- Consistent with SDK patterns

**Structure:**
```go
func (c *Client) Get(ctx context.Context, path string, out interface{}) error
func (c *Client) Post(ctx context.Context, path string, body interface{}, out interface{}) error
func (c *Client) Put(ctx context.Context, path string, body interface{}, out interface{}) error
func (c *Client) Delete(ctx context.Context, path string, out interface{}) error
```

### 6. WithOptions Method
**Decision:** Implement client cloning with modified options

**Rationale:**
- Enables creating derived clients with different settings
- Non-mutating (functional pattern)
- Useful for per-request configuration

**Structure:**
```go
func (c *Client) WithOptions(opts ...Option) *Client
```

## Risks / Trade-offs

**[API complexity]** → Many optional parameters may confuse users
- **Mitigation:** Clear documentation; mark advanced parameters as such

**[Raw response overhead]** → Storing raw body increases memory usage
- **Mitigation:** Only store when explicitly requested; document memory implications

**[Reasoning model availability]** → Not all models support reasoning
- **Mitigation:** Document which models support reasoning features
