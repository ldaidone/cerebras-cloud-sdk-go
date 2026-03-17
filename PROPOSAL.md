# Proposal: Cerebras Cloud Go SDK (Pure Go Port)

## 1. Executive Summary

The goal of this project is to provide a high-performance, idiomatic Go alternative to the official Python SDK. This port will focus on strict typing, minimal dependencies, and native Go concurrency patterns to support low-latency inference streaming.

## 2. Tech Stack & Strategy

* **Language:** Go 1.24+ (utilizing `context` and `slog` for structured logging).
* **HTTP Layer:** `net/http` (Standard Library) to ensure zero-dependency core.
* **Architecture:** Spec-Driven Development (SDD). We will define the API contract in a local specification file to auto-generate or manually validate structs.
* **Concurrency:** Channels and `context.Context` for managing streaming timeouts and resource cleanup.

## 3. Directory Structure

```markdown
.
├── cmd/
│   └── cerebras-cli/      # Example CLI tool
├── internal/
│   ├── transport/         # HTTP client overrides and retries
│   └── stream/            # SSE (Server-Sent Events) decoder logic
├── pkg/
│   ├── cerebras/
│   │   ├── client.go      # Primary SDK entry point
│   │   ├── options.go     # Functional options (WithTimeout, WithAPIKey)
│   │   ├── chat.go        # Chat Completion methods
│   │   └── types.go       # Structs derived from API Spec
└── spec/
    └── openapi.yaml       # The source of truth for the SDK
```

## 4. Implementation Strategy

### Phase 1: The "Source of Truth"

We will start by capturing the Cerebras OpenAI-compatible spec. This allows us to map `pydantic` models from Python directly into Go `structs` with precise JSON tags.

### Phase 2: The Core Client

The client will use **Functional Options** for configuration.

- `NewClient(opts ...Option)`

- This prevents "breaking changes" when adding new configuration parameters (like custom base URLs for local proxies or private endpoints).

### Phase 3: Streaming & Channels

Unlike Python’s `yield`, the Go SDK will implement a `Stream` object that provides a `Recv()` method. Internally, this will use a `bufio.Scanner` to minimize memory allocations during high-throughput token streaming.

## 5. Build & Development Steps

1. **Initialize:** `go mod init github.com/youruser/cerebras-go`

2. **Define Types:** Translate the chat completion request/response bodies into `types.go`.

3. **Implement Transport:** Build the authenticated `http.Client` wrapper.

4. **Streaming Logic:** Build the SSE parser in `internal/stream`.

5. **Validation:** Create a test suite that mocks the Cerebras API responses to ensure parity with the Python SDK behavior.

## 6. Key Considerations

- **Performance:** Avoid `interface{}` where possible to keep reflection costs low during JSON unmarshaling.

- **Error Handling:** Implement custom error types (e.g., `APIError`, `RateLimitError`) that allow users to programmatically handle 429s or 500s.

- **Local AI Parity:** Ensure the `BaseURL` can be overridden to support local LLM proxies (Ollama, LM Studio) during development.
