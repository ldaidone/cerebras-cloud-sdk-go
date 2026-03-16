## Why

The Go SDK lacks SSE (Server-Sent Events) streaming support, which is essential for real-time token streaming from LLM responses. This change adds streaming capabilities enabling users to receive model responses as they are generated, improving perceived latency and enabling interactive applications.

## What Changes

- Implement SSE streaming client for chat completions
- Add `POST /v1/chat/completions` endpoint with `stream=true` support
- Implement streaming response parser for SSE format
- Add `StreamResponse` iterator/channel-based API for consuming chunks
- Support `stream_options` with `include_usage` parameter
- Integrate with existing types and error handling
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `sse-streaming-core`: SSE streaming infrastructure including HTTP streaming, event parsing, and response iterator
- `chat-completions-streaming`: Streaming support for chat completions endpoint with CreateStream method

### Modified Capabilities
- None

## Impact

- **Affected code**: New streaming package in `pkg/cerebras/streaming.go`
- **Dependencies**: Requires types from `types-definitions` change
- **APIs**: New public API `client.ChatCompletions.CreateStream(ctx, request)` returning stream iterator
- **Systems**: Enables real-time response streaming for chat completions
