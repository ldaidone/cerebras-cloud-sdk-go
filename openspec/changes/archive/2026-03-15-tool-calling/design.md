## Context

The Go SDK has completed its foundation layer, type definitions, chat completions, and streaming. Tool/function calling is an advanced feature that enables models to invoke external functions. The Python SDK provides tool definitions and tool choice parameters. This change adds equivalent functionality to the Go SDK.

Key constraints:
- API compatibility with Python SDK tool calling patterns
- JSON Schema for function parameters (stdlib encoding/json)
- Integration with existing ChatCompletionRequest types
- Support for parallel tool calls

## Goals / Non-Goals

**Goals:**
- Define Tool, Function, and ToolParameters types
- Define ToolChoice type (string "none"/"auto"/"required" or object)
- Define ToolCall and FunctionCall for responses
- Integrate tools, tool_choice, parallel_tool_calls with chat completions
- Provide helper functions for tool definition
- Comprehensive examples and documentation

**Non-Goals:**
- Automatic function execution (user handles function calls)
- Function response submission (separate API call pattern)
- Streaming tool calls (covered in sse-streaming)

## Decisions

### 1. Tool Definition Structure
**Decision:** Define Tool struct with type="function" and Function field

**Rationale:**
- Matches OpenAI/Cerebras API format
- Allows future extension with other tool types
- Clear separation between tool metadata and function definition

**Structure:**
```go
type Tool struct {
    Type     string   `json:"type"`  // "function"
    Function Function `json:"function"`
}

type Function struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    Parameters  map[string]interface{} `json:"parameters"`  // JSON Schema
}
```

### 2. JSON Schema for Parameters
**Decision:** Use map[string]interface{} for function parameters

**Rationale:**
- JSON Schema is flexible and dynamic
- Avoids external JSON Schema library dependencies
- Users can define any schema structure
- Consistent with Python SDK's dict approach

**Alternatives considered:**
- Strongly typed JSON Schema struct: Too verbose, limits flexibility
- External JSON Schema library: Adds dependency

### 3. Tool Choice Type
**Decision:** Use interface{} for ToolChoice (string or object)

**Rationale:**
- API accepts string ("none", "auto", "required") or object ({"type": "function", "function": {...}})
- interface{} provides flexibility
- Custom JSON marshaler can handle both cases

**Implementation:**
```go
type ToolChoice interface{}

// String choice
toolChoice := "auto"

// Object choice
toolChoice := ToolChoiceFunction{
    Type: "function",
    Function: ToolChoiceFunctionName{Name: "my_function"},
}
```

### 4. Tool Call in Response
**Decision:** Add ToolCalls array to Message struct

**Rationale:**
- Model may call multiple tools in single response
- Consistent with Python SDK structure
- Each tool call has unique ID for tracking

**Structure:**
```go
type ToolCall struct {
    ID       string       `json:"id"`
    Type     string       `json:"type"`
    Function FunctionCall `json:"function"`
}

type FunctionCall struct {
    Name      string `json:"name"`
    Arguments string `json:"arguments"`  // JSON string
}
```

### 5. Function Arguments as JSON String
**Decision:** Store function arguments as JSON string (not parsed object)

**Rationale:**
- API returns arguments as JSON string
- User can parse into desired type
- Avoids assumptions about argument structure
- Consistent with Python SDK

### 6. Helper Functions
**Decision:** Provide helper functions for common tool definition patterns

**Rationale:**
- Simplifies tool definition for users
- Reduces boilerplate
- Provides examples of JSON Schema construction

**Helpers:**
```go
func DefineFunction(name, description string, params map[string]interface{}) Tool
func WithRequiredParams(params map[string]interface{}, required []string) map[string]interface{}
```

## Risks / Trade-offs

**[JSON Schema complexity]** → Users must understand JSON Schema format
- **Mitigation:** Provide clear examples and helper functions

**[Type safety]** → map[string]interface{} lacks compile-time type checking
- **Mitigation:** Document expected schema format; provide examples

**[Function execution]** → Users must manually execute functions
- **Mitigation:** Provide example pattern for function execution loop
