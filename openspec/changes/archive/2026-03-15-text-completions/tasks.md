## 1. Text Completion Types

- [x] 1.1 Define TextCompletionRequest struct with model, prompt, and optional fields
- [x] 1.2 Define TextCompletion struct with id, object, created, model, choices, usage
- [x] 1.3 Define TextCompletionChoice struct with text, index, finish_reason
- [x] 1.4 Define TextCompletionStreamResponse struct for streaming
- [x] 1.5 Add JSON tags to all text completion types

## 2. Text Completions Service

- [x] 2.1 Create TextCompletionsService struct in pkg/cerebras/text_completions.go
- [x] 2.2 Add textCompletions field to Client struct
- [x] 2.3 Initialize TextCompletionsService in NewClient or constructor

## 3. Create Method Implementation

- [x] 3.1 Implement Create method signature with context and request
- [x] 3.2 Implement HTTP POST /v1/completions request
- [x] 3.3 Handle response deserialization to TextCompletion struct
- [x] 3.4 Add parameter validation (model required, prompt required)

## 4. Streaming Support

- [x] 4.1 Implement CreateStream method for streaming text completions
- [x] 4.2 Handle SSE parsing for text completion chunks
- [x] 4.3 Support stream_options with include_usage

## 5. Error Handling

- [x] 5.1 Integrate with existing error hierarchy
- [x] 5.2 Handle validation errors
- [x] 5.3 Handle API errors during streaming

## 6. Documentation

- [x] 6.1 Add GoDoc comments to text completion types and methods
- [x] 6.2 Create example showing basic text completion
- [x] 6.3 Create example showing streaming text completion
- [x] 6.4 Document use cases (legacy endpoint, non-chat scenarios)
