## Why

The Go SDK lacks advanced features present in the Python SDK including service tiers, reasoning model support, and raw response access. This change adds these advanced capabilities enabling users to access priority processing, reasoning models with chain-of-thought, and low-level HTTP response inspection for debugging and advanced use cases.

## What Changes

- Add service_tier parameter support (auto, flex, priority)
- Add reasoning_effort and reasoning_format parameters for reasoning models
- Add clear_thinking and disable_reasoning parameters
- Implement raw response access pattern for HTTP inspection
- Add HTTP verb helpers (Get, Post, Put, Delete) for custom API calls
- Add WithOptions method for client configuration chaining
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `service-tiers`: Service tier parameter support for priority processing control
- `reasoning-support`: Reasoning model parameters (reasoning_effort, reasoning_format, clear_thinking, disable_reasoning)
- `raw-response-access`: Raw HTTP response access pattern for debugging and inspection
- `http-verbs`: HTTP verb helpers for custom API calls

### Modified Capabilities
- None

## Impact

- **Affected code**: Extended types in `pkg/cerebras/`, new utility methods
- **Dependencies**: Requires types from `types-definitions` change
- **APIs**: Extended ChatCompletionRequest with advanced parameters, new client methods
- **Systems**: Enables advanced use cases and debugging capabilities
