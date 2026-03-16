## Context

The Go SDK has completed its foundation layer, type definitions, chat completions, streaming, tool calling, and response formats. Text Completions is a legacy API endpoint that provides simple text generation without the chat message format. While chat completions is preferred for most use cases, text completions is still supported and useful for simple completion tasks.

Key constraints:
- API compatibility with Python SDK text completions patterns
- Similar structure to chat completions but with prompt instead of messages
- Support for both non-streaming and streaming modes
- Integration with existing transport layer

## Goals / Non-Goals

**Goals:**
- Implement TextCompletionsService with Create and CreateStream methods
- Support POST /v1/completions endpoint
- Define TextCompletionRequest and TextCompletion types
- Support all standard parameters (model, prompt, max_tokens, temperature, etc.)
- Support streaming for text completions
- Proper error handling with typed errors
- Comprehensive examples and documentation

**Non-Goals:**
- Advanced text completion features not supported by API
- Batch completions (API doesn't support)
- Fine-tuning integration (separate API)

## Decisions

### 1. Service Structure
**Decision:** Implement as TextCompletionsService struct with method receiver

**Rationale:**
- Matches Python SDK's `client.completions` namespace
- Groups related operations logically
- Consistent with ChatCompletionsService pattern

**Structure:**
```go
type TextCompletionsService struct {
    client *Client
}

func (s *TextCompletionsService) Create(ctx context.Context, request *TextCompletionRequest) (*TextCompletion, error)
func (s *TextCompletionsService) CreateStream(ctx context.Context, request *TextCompletionRequest) (<-chan TextCompletionStreamResponse, <-chan error, error)
```

### 2. Request Type Structure
**Decision:** Define TextCompletionRequest with prompt (string or []string)

**Rationale:**
- API accepts single prompt or array of prompts
- Use interface{} or separate fields for flexibility
- Consistent with Python SDK's Union[str, List[str]]

**Structure:**
```go
type TextCompletionRequest struct {
    Model            string    `json:"model"`
    Prompt           string    `json:"prompt"`  // or []string for batch
    MaxTokens        *int      `json:"max_tokens,omitempty"`
    Temperature      *float64  `json:"temperature,omitempty"`
    // ... other parameters
}
```

### 3. Response Type Structure
**Decision:** Define TextCompletion with choices array containing text

**Rationale:**
- Similar to ChatCompletion but with text instead of message
- Each choice contains text, finish_reason, and index
- Consistent with API response format

**Structure:**
```go
type TextCompletion struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   *Usage   `json:"usage,omitempty"`
}

type Choice struct {
    Text         string  `json:"text"`
    Index        int     `json:"index"`
    FinishReason string  `json:"finish_reason"`
}
```

### 4. Streaming Response Type
**Decision:** Define TextCompletionStreamResponse similar to chat streaming

**Rationale:**
- Consistent with chat completions streaming pattern
- Delta contains text instead of message delta
- Reuses SSE streaming infrastructure

### 5. Parameter Handling
**Decision:** Support same optional parameters as chat completions where applicable

**Rationale:**
- API supports similar parameters (temperature, max_tokens, etc.)
- Some chat-specific parameters don't apply (messages, tools)
- Document parameter differences from chat completions

## Risks / Trade-offs

**[Legacy API]** → Text completions may be deprecated in future
- **Mitigation:** Document that chat completions is preferred; maintain compatibility

**[Prompt vs messages]** → Users may confuse text and chat completions
- **Mitigation:** Clear documentation on when to use each; examples

**[Batch prompts]** → API supports array of prompts
- **Mitigation:** Start with single prompt; add batch support in future update
