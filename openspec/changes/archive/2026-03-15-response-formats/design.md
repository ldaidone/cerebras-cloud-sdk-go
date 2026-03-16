## Context

The Go SDK has completed its foundation layer, type definitions, chat completions, streaming, and tool calling. Response format controls are an important feature for structured data extraction, enabling models to output valid JSON. The Python SDK provides response_format parameter with json_object and json_schema options.

Key constraints:
- API compatibility with Python SDK response_format patterns
- JSON Schema for schema-constrained output
- Integration with existing ChatCompletionRequest types
- Support for both json_object and json_schema modes

## Goals / Non-Goals

**Goals:**
- Define ResponseFormat type with type field
- Define JSONSchema type for schema-constrained output
- Integrate response_format with chat completions request
- Provide helper functions for common response format patterns
- Comprehensive examples and documentation

**Non-Goals:**
- Automatic JSON parsing/validation (user handles parsing)
- Schema validation of responses (user validates)
- Response format for streaming (API limitation)

## Decisions

### 1. ResponseFormat Structure
**Decision:** Define ResponseFormat as struct with type and optional json_schema fields

**Rationale:**
- Matches OpenAI/Cerebras API format
- Clear separation between format type and schema definition
- Extensible for future format types

**Structure:**
```go
type ResponseFormatType string
const (
    ResponseFormatText       ResponseFormatType = "text"
    ResponseFormatJSONObject ResponseFormatType = "json_object"
    ResponseFormatJSONSchema ResponseFormatType = "json_schema"
)

type ResponseFormat struct {
    Type       ResponseFormatType `json:"type"`
    JSONSchema *JSONSchema        `json:"json_schema,omitempty"`
}
```

### 2. JSON Schema Structure
**Decision:** Define JSONSchema with name, description, schema, and strict fields

**Rationale:**
- Matches API specification
- name and schema are required for json_schema mode
- strict enables strict schema adherence (optional)

**Structure:**
```go
type JSONSchema struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    Schema      map[string]interface{} `json:"schema"`
    Strict      *bool                  `json:"strict,omitempty"`
}
```

### 3. Schema as map[string]interface{}
**Decision:** Use map[string]interface{} for JSON schema definition

**Rationale:**
- JSON Schema is flexible and dynamic
- Avoids external JSON Schema library dependencies
- Users can define any schema structure
- Consistent with Python SDK's dict approach

### 4. Integration with Chat Completions
**Decision:** Add ResponseFormat field to ChatCompletionRequest

**Rationale:**
- Consistent with other optional parameters
- Pointer allows omitting when not specified
- Clear API for users

**Implementation:**
```go
type ChatCompletionRequest struct {
    // ... other fields
    ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}
```

### 5. Helper Functions
**Decision:** Provide helper functions for common response format patterns

**Rationale:**
- Simplifies common use cases
- Reduces boilerplate
- Provides examples of proper usage

**Helpers:**
```go
func ResponseFormatJSON() *ResponseFormat
func ResponseFormatJSONWithSchema(name string, schema map[string]interface{}) *ResponseFormat
```

## Risks / Trade-offs

**[JSON Schema complexity]** → Users must understand JSON Schema format
- **Mitigation:** Provide clear examples and helper functions

**[Type safety]** → map[string]interface{} lacks compile-time type checking
- **Mitigation:** Document expected schema format; provide examples

**[JSON validation]** → Model may not always produce valid JSON
- **Mitigation:** Document that users should handle parse errors
