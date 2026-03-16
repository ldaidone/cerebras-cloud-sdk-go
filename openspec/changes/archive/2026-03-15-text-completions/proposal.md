## Why

The Go SDK lacks the Text Completions API endpoint, which is a legacy but still-supported interface for simple text generation tasks. This change adds the text completions endpoint enabling users to generate text completions for prompts without the chat message format.

## What Changes

- Implement `TextCompletionsService` with `Create` method
- Add `POST /v1/completions` endpoint support
- Define `TextCompletionRequest` and `TextCompletion` response types
- Support all standard parameters (model, prompt, max_tokens, temperature, etc.)
- Add streaming support for text completions
- Integrate with existing HTTP transport and error handling
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `text-completions-endpoint`: Non-streaming text completions API implementation with full parameter support
- `text-completions-streaming`: Streaming support for text completions endpoint

### Modified Capabilities
- None

## Impact

- **Affected code**: New service in `pkg/cerebras/text_completions.go`
- **Dependencies**: Requires types from `types-definitions` and `sse-streaming` changes
- **APIs**: New public API `client.TextCompletions.Create(ctx, request)`
- **Systems**: Enables legacy text completion workflows
