## Context

The Go SDK has completed its foundation layer (client infrastructure, error handling, HTTP transport) but lacks type definitions for API entities. The Python SDK provides comprehensive type definitions using Pydantic models. This change translates those patterns to idiomatic Go structs with proper JSON tags and validation.

Key constraints:
- Must maintain API compatibility with Python SDK field names and behavior
- Go-specific patterns (pointers for optional fields, time.Time for timestamps)
- Zero external dependencies (stdlib only)
- Support for JSON marshaling/unmarshaling with proper null handling

## Goals / Non-Goals

**Goals:**
- Define all core types for Chat Completions API (Message, ChatCompletionRequest, ChatCompletion, Choice, Usage, etc.)
- Define types for Models API (Model, ModelList)
- Define types for streaming (StreamResponse, Delta, StreamOptions)
- Define types for tool calling (Tool, ToolChoice, Function)
- Define types for response formats (ResponseFormat, JSONSchema)
- Define model constants for all supported models
- Ensure JSON serialization compatibility with Cerebras Cloud API
- Provide GoDoc documentation for all exported types

**Non-Goals:**
- Implementation of API endpoints (covered in separate changes)
- Validation logic beyond JSON tags (handled by API server)
- Builder patterns or helper methods (can be added later)

## Decisions

### 1. Package Organization
**Decision:** Place all types in `pkg/cerebras/` package directly (not subpackage)

**Rationale:**
- Follows Go convention for SDK packages (e.g., `github.com/aws/aws-sdk-go-v2/service/s3/types`)
- Reduces import path complexity for users
- Keeps related types together with client code

**Alternatives considered:**
- `pkg/cerebras/types/` subpackage: Adds verbosity (`types.Message` vs `cerebras.Message`)
- Separate `types` package: Creates circular dependency risk with client

### 2. Optional Fields as Pointers
**Decision:** Use pointers for optional fields (`*string`, `*int`, `*float64`)

**Rationale:**
- Distinguishes between "not set" (nil) and "zero value" (0, "", false)
- Matches Python SDK's Optional[] pattern
- Enables proper JSON omitempty behavior

**Example:**
```go
type ChatCompletionRequest struct {
    Model       string  `json:"model"`
    Temperature *float64 `json:"temperature,omitempty"`
    MaxTokens   *int     `json:"max_tokens,omitempty"`
}
```

### 3. Enum Types as String Aliases
**Decision:** Define enums as string type aliases with constants

**Rationale:**
- Type safety without external dependencies
- Easy JSON marshaling via string conversion
- Clear IDE autocomplete support

**Example:**
```go
type MessageRole string
const (
    MessageRoleSystem    MessageRole = "system"
    MessageRoleUser      MessageRole = "user"
    MessageRoleAssistant MessageRole = "assistant"
)
```

### 4. Timestamp Handling
**Decision:** Use `int64` for Unix timestamps (not `time.Time`)

**Rationale:**
- API returns Unix timestamps as integers
- Avoids unnecessary time parsing/marshaling complexity
- Users can convert to `time.Time` via `time.Unix()` when needed

### 5. Model Constants
**Decision:** Define model constants as string constants in main package

**Rationale:**
- Prevents typos in model names
- Provides IDE autocomplete
- Documents supported models in one place

## Risks / Trade-offs

**[Type verbosity]** → Go types are more verbose than Python's Pydantic models
- **Mitigation:** Focus on essential fields first; add convenience methods later

**[Pointer complexity]** → Users must work with pointers for optional fields
- **Mitigation:** Provide helper functions (e.g., `PtrString()`, `PtrInt()`) in future update

**[API drift]** → Python SDK may add new fields
- **Mitigation:** Design types to be extensible; add unknown fields handling if needed

**[Breaking changes]** → Future API changes may require type modifications
- **Mitigation:** Use pointers for new optional fields to maintain backward compatibility
