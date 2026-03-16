## Why

The Go SDK currently lacks type definitions for core API entities, making it impossible to use type-safe requests and responses when interacting with Cerebras Cloud APIs. This change introduces comprehensive type definitions matching the Python SDK's API compatibility.

## What Changes

- Add core type definitions for Chat Completions API (Message, ChatCompletionRequest, ChatCompletion, etc.)
- Add type definitions for Models API (Model, ModelList)
- Add type definitions for streaming responses (StreamResponse, Delta, StreamOptions)
- Add type definitions for tool calling (Tool, ToolChoice, Function)
- Add type definitions for response formats (ResponseFormat, JSONSchema)
- Add model constants (Llama31_8b, Llama31_70b, etc.)
- Add type definitions for text completions endpoint
- **BREAKING**: Existing code using `map[string]interface{}` for API payloads will need to migrate to typed structs

## Capabilities

### New Capabilities
- `type-definitions`: Core type definitions for all API entities including messages, completions, models, and common parameters
- `model-constants`: Predefined constants for supported model identifiers

### Modified Capabilities
- None

## Impact

- **Affected code**: New package `pkg/cerebras/types` or types in `pkg/cerebras/`
- **Dependencies**: None (stdlib only)
- **Systems**: Foundation for all subsequent API endpoint implementations
- **APIs**: Enables Chat Completions, Models, Text Completions, Streaming, and Tool Calling implementations
