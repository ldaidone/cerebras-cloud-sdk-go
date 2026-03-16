## Context

The Go SDK has completed its foundation layer, type definitions, and non-streaming chat completions. SSE streaming is a critical feature for LLM APIs, enabling real-time token streaming. The Python SDK uses async generators for streaming; Go will use channels and iterators following Go concurrency patterns.

Key constraints:
- API compatibility with Python SDK streaming behavior
- SSE format parsing (data: JSON lines)
- Go-specific patterns (channels, context cancellation)
- Integration with existing transport layer
- Proper resource cleanup (response body closing)

## Goals / Non-Goals

**Goals:**
- Implement SSE streaming client for HTTP streaming
- Parse SSE event format (data: JSON lines)
- Implement CreateStream method for chat completions streaming
- Return channel-based or iterator-based stream response
- Support stream_options with include_usage
- Proper error handling and context cancellation
- Comprehensive examples and documentation

**Non-Goals:**
- Streaming for text completions (covered in text-completions change)
- Advanced streaming features (tool calling streams, etc.)
- WebSocket or other streaming protocols

## Decisions

### 1. Streaming API Pattern
**Decision:** Use channel-based streaming with StreamResponse struct

**Rationale:**
- Idiomatic Go pattern for streaming data
- Enables range loops and select statements
- Integrates with Go's concurrency model
- Clear semantics for stream lifecycle

**Structure:**
```go
type StreamingChatCompletionsService struct {
    client *Client
}

func (s *StreamingChatCompletionsService) Create(ctx context.Context, request *ChatCompletionRequest) (<-chan StreamResponse, error)
```

### 2. Stream Response Channel
**Decision:** Return receive-only channel of StreamResponse

**Rationale:**
- Channel naturally represents stream of events
- Receive-only prevents user from sending on channel
- Channel closure signals end of stream
- Enables non-blocking consumption with select

**Implementation:**
```go
stream, err := client.ChatCompletions.CreateStream(ctx, request)
for response := range stream {
    // Process chunk
}
```

### 3. SSE Parsing
**Decision:** Implement custom SSE parser for data: lines

**Rationale:**
- SSE format is simple (data: JSON\n\n)
- No external dependencies needed
- Full control over error handling
- Can handle edge cases (multi-line data, etc.)

**Format:**
```
data: {"id":"...","choices":[...]}

data: {"id":"...","choices":[...],"usage":{}}

data: [DONE]
```

### 4. HTTP Client Configuration
**Decision:** Use separate HTTP client for streaming with no timeout

**Rationale:**
- Streaming requests may take extended time
- Default timeout may interrupt long streams
- User can still cancel via context
- Separate client avoids affecting non-streaming requests

### 5. Error Handling in Stream
**Decision:** Send errors through channel and return via separate error channel

**Rationale:**
- Distinguishes between stream errors and iteration completion
- Enables proper error handling during iteration
- Consistent with Go error handling patterns

**Implementation:**
```go
func (s *StreamingChatCompletionsService) Create(ctx context.Context, request *ChatCompletionRequest) (<-chan StreamResponse, <-chan error, error)
```

### 6. Resource Cleanup
**Decision:** Use context cancellation and automatic body closing

**Rationale:**
- Prevents resource leaks
- User doesn't need to manually close stream
- Context cancellation propagates to HTTP request
- Defer cleanup in goroutine

## Risks / Trade-offs

**[Channel buffering]** → Unbuffered channels may block
- **Mitigation:** Use small buffer (e.g., 10) to prevent blocking

**[Error handling complexity]** → Two-channel pattern (data + errors) is complex
- **Mitigation:** Provide clear examples and helper functions

**[Memory usage]** → Long streams may accumulate data
- **Mitigation:** Document that users should process chunks as received

**[Connection pooling]** → Long-lived streaming connections affect pooling
- **Mitigation:** Use separate HTTP client for streaming
