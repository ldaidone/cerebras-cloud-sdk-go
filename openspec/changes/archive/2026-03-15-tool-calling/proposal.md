## Why

The Go SDK lacks tool/function calling support, which enables LLM models to invoke external functions and APIs. This change adds tool calling capabilities allowing users to define functions that models can call, enabling complex workflows and integrations with external systems.

## What Changes

- Add tool definition types (Tool, Function, ToolParameters)
- Add tool choice types (ToolChoice, ToolChoiceFunction)
- Add tool call types in response (ToolCall, FunctionCall)
- Implement tools and tool_choice parameters in chat completions
- Support parallel_tool_calls parameter
- Add helper functions for tool definition
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `tool-definitions`: Type definitions for tools, functions, and tool choice
- `tool-calling-support`: Tool calling integration with chat completions endpoint

### Modified Capabilities
- None

## Impact

- **Affected code**: New types in `pkg/cerebras/`, chat completions integration
- **Dependencies**: Requires types from `types-definitions` change
- **APIs**: Extended ChatCompletionRequest with tools, tool_choice, parallel_tool_calls
- **Systems**: Enables function calling workflows and external integrations
