## Context

The Go SDK has completed its foundation layer (client infrastructure, error handling, HTTP transport) and type definitions are being added. The Chat Completions API is the primary interface for LLM interactions, matching the Python SDK's `client.chat.completions.create()` pattern.

Key constraints:
- API compatibility with Python SDK method signatures and behavior
- Support for all parameters defined in Cerebras Cloud API
- Integration with existing transport layer (retry logic, timeouts, context propagation)
- Non-streaming implementation first (streaming in separate change)

## Goals / Non-Goals

**Goals:**
- Implement `ChatCompletionsService` with `Create` method
- Support POST /v1/chat/completions endpoint
- Handle all standard parameters (model, messages, temperature, max_tokens, top_p, stop, n, etc.)
- Return typed ChatCompletion response
- Proper error handling with typed errors
- Comprehensive examples and documentation

**Non-Goals:**
- Streaming support (covered in sse-streaming change)
- Tool calling support (covered in tool-calling change)
- Response format controls (covered in response-formats change)
- Service tiers and advanced parameters (covered in advanced-features change)

## Decisions

### 1. Service Structure
**Decision:** Implement as `ChatCompletionsService` struct with method receiver

**Rationale:**
- Matches Python SDK's `client.chat.completions` namespace
- Groups related operations logically
- Enables future extension with additional methods

**Structure:**
```go
type ChatCompletionsService struct {
    client *Client
}

func (s *ChatCompletionsService) Create(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletion, error)
```

### 2. Method Signature
**Decision:** Use functional options pattern for optional parameters

**Rationale:**
- Consistent with Go best practices for APIs with many optional parameters
- Type-safe alternative to variadic map[string]interface{}
- Enables IDE autocomplete for options

**Alternatives considered:**
- Direct struct parameter: Requires users to set pointers manually
- Builder pattern: More verbose, adds complexity

**Implementation:**
```go
func (s *ChatCompletionsService) Create(ctx context.Context, model string, messages []Message, opts ...ChatCompletionOption) (*ChatCompletion, error)
```

### 3. Request/Response Handling
**Decision:** Use existing transport layer's `sendRequest` method

**Rationale:**
- Reuses retry logic, error handling, and JSON marshaling
- Maintains consistency across all endpoints
- Leverages existing context propagation and timeout handling

### 4. Parameter Validation
**Decision:** Minimal client-side validation (defer to API server)

**Rationale:**
- API server is source of truth for validation rules
- Reduces SDK maintenance burden when API changes
- Avoids false positives from client-side validation

**Validation performed:**
- Required fields present (model, messages)
- Messages array not empty

### 5. Error Handling
**Decision:** Return typed errors from existing error hierarchy

**Rationale:**
- Consistent with SDK's error handling approach
- Enables errors.Is/as pattern matching
- Provides clear error messages with HTTP status codes

## Risks / Trade-offs

**[API parameter drift]** → Cerebras API may add new parameters
- **Mitigation:** Design ChatCompletionRequest to be extensible; add unknown fields handling

**[Functional options complexity]** → Users unfamiliar with functional options pattern
- **Mitigation:** Provide clear examples and documentation

**[Large response handling]** → Chat completions may return large responses
- **Mitigation:** Rely on existing HTTP client configuration; document memory considerations
