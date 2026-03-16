## 1. Message Types

- [x] 1.1 Define MessageRole enum type with constants (system, user, assistant)
- [x] 1.2 Define Message struct with role, content, name, and optional fields
- [x] 1.3 Add JSON tags to Message struct for proper serialization

## 2. Chat Completion Request Types

- [x] 2.1 Define ChatCompletionRequest struct with required fields (model, messages)
- [x] 2.2 Add optional fields as pointers (temperature, max_tokens, top_p, stop, n, etc.)
- [x] 2.3 Add JSON tags with snake_case naming
- [x] 2.4 Add common parameters (frequency_penalty, presence_penalty, logprobs, etc.)

## 3. Chat Completion Response Types

- [x] 3.1 Define Choice struct with message, finish_reason, and index
- [x] 3.2 Define ChatCompletion struct with id, object, created, model, choices, usage
- [x] 3.3 Define Usage struct with token counts and optional details
- [x] 3.4 Define TimeInfo struct with timing metrics
- [x] 3.5 Add FinishReason enum type with constants

## 4. Model Types

- [x] 4.1 Define Model struct with id, object, created, owned_by
- [x] 4.2 Define ModelList struct with object and data array

## 5. Streaming Types

- [x] 5.1 Define Delta struct with role, content, and optional fields
- [x] 5.2 Define StreamResponse struct with choices containing deltas
- [x] 5.3 Define StreamOptions struct with include_usage field

## 6. Tool Calling Types

- [x] 6.1 Define Function struct with name, description, parameters
- [x] 6.2 Define Tool struct with type and function fields
- [x] 6.3 Define ToolChoice type (string or object)
- [x] 6.4 Define ToolCall struct with id, type, and function fields

## 7. Response Format Types

- [x] 7.1 Define ResponseFormat struct with type field
- [x] 7.2 Define JSONSchema struct with name, description, schema
- [x] 7.3 Add ResponseFormat union type for json_object and json_schema

## 8. Model Constants

- [x] 8.1 Define model identifier constants (Llama31_8b, Llama31_70b, etc.)
- [x] 8.2 Add GoDoc comments for each constant

## 9. Error Response Types

- [x] 9.1 Define ErrorResponse struct matching API error format
- [x] 9.2 Define Error struct with message, type, param, code fields

## 10. Helper Functions

- [x] 10.1 Add pointer helper functions (PtrString, PtrInt, PtrFloat64)
- [x] 10.2 Add example usage in doc comments
