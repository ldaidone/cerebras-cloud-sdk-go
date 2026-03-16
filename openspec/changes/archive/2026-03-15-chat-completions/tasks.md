## 1. Service Structure

- [x] 1.1 Create ChatCompletionsService struct in pkg/cerebras/chat_completions.go
- [x] 1.2 Add chatCompletions field to Client struct
- [x] 1.3 Initialize ChatCompletionsService in NewClient or constructor

## 2. Create Method Implementation

- [x] 2.1 Implement Create method signature with context, model, messages, and options
- [x] 2.2 Implement ChatCompletionOption functional option type
- [x] 2.3 Add functional option helpers (WithTemperature, WithMaxTokens, etc.)
- [x] 2.4 Build ChatCompletionRequest from parameters and options

## 3. HTTP Request Handling

- [x] 3.1 Implement POST /v1/chat/completions request using transport layer
- [x] 3.2 Set proper request headers (Authorization, Content-Type, User-Agent)
- [x] 3.3 Serialize request body to JSON with snake_case field names
- [x] 3.4 Handle response deserialization to ChatCompletion struct

## 4. Parameter Validation

- [x] 4.1 Add validation for required model parameter
- [x] 4.2 Add validation for required messages parameter
- [x] 4.3 Add validation for non-empty messages array

## 5. Error Handling

- [x] 5.1 Integrate with existing error hierarchy for API errors
- [x] 5.2 Handle 400 Bad Request errors
- [x] 5.3 Handle 401 Authentication errors
- [x] 5.4 Handle 429 Rate Limit errors
- [x] 5.5 Handle 5xx Server errors

## 6. Context and Retry Support

- [x] 6.1 Ensure context propagation for cancellation
- [x] 6.2 Verify retry logic applies to chat completions requests
- [x] 6.3 Test context timeout behavior

## 7. Documentation

- [x] 7.1 Add GoDoc comments to ChatCompletionsService and Create method
- [x] 7.2 Add example usage in doc comments
- [x] 7.3 Create example file showing basic chat completion usage
