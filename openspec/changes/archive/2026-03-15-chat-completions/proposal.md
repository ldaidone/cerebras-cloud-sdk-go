## Why

The Go SDK lacks the Chat Completions API endpoint implementation, which is the primary interface for interacting with Cerebras Cloud's LLM models. This change adds the non-streaming chat completions endpoint enabling users to send messages and receive model responses.

## What Changes

- Implement `ChatCompletions` service with `Create` method
- Add `POST /v1/chat/completions` endpoint support
- Support all standard parameters (model, messages, temperature, max_tokens, etc.)
- Return typed `ChatCompletion` response with choices, usage, and metadata
- Integrate with existing HTTP transport and error handling
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `chat-completions-endpoint`: Non-streaming chat completions API implementation with full parameter support

### Modified Capabilities
- None

## Impact

- **Affected code**: New service in `pkg/cerebras/chat_completions.go`
- **Dependencies**: Requires types from `types-definitions` change
- **APIs**: New public API `client.ChatCompletions.Create(ctx, request)`
- **Systems**: Foundation for streaming and advanced chat features
