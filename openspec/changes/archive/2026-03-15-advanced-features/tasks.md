## 1. Service Tier Types

- [x] 1.1 Define ServiceTier enum with constants (auto, default, flex, priority)
- [x] 1.2 Add ServiceTier field to ChatCompletionRequest

## 2. Reasoning Support

- [x] 2.1 Define ReasoningEffort enum (low, medium, high)
- [x] 2.2 Define ReasoningFormat enum (none, parsed, text_parsed, raw, hidden)
- [x] 2.3 Add reasoning_effort, reasoning_format to ChatCompletionRequest
- [x] 2.4 Add clear_thinking, disable_reasoning to ChatCompletionRequest

## 3. Raw Response Access

- [x] 3.1 Define APIResponse struct with StatusCode, Headers, Body
- [x] 3.2 Implement WithRawResponse method for ChatCompletionsService
- [x] 3.3 Document debugging use cases

## 4. Documentation

- [x] 4.1 Add GoDoc comments to advanced feature types
- [x] 4.2 Document service tier trade-offs
- [x] 4.3 Document reasoning model capabilities
- [x] 4.4 Document raw response access patterns
