## ADDED Requirements

### Requirement: Chat completions streaming service
The system SHALL provide a StreamingChatCompletionsService accessible via client.ChatCompletions with a CreateStream method for streaming chat completions.

#### Scenario: Service access
- **WHEN** a user initializes the client
- **THEN** client.ChatCompletions provides access to streaming service via CreateStream method

#### Scenario: CreateStream method signature
- **WHEN** a user calls CreateStream method
- **THEN** method accepts context, ChatCompletionRequest, and returns stream channel and error

### Requirement: POST /v1/chat/completions streaming endpoint
The system SHALL implement HTTP POST to /v1/chat/completions with stream=true parameter.

#### Scenario: Request endpoint
- **WHEN** CreateStream method is called
- **THEN** HTTP POST request is sent to /v1/chat/completions

#### Scenario: Stream parameter
- **WHEN** request is sent
- **THEN** request body includes stream=true

#### Scenario: Request headers
- **WHEN** request is sent
- **THEN** headers include Authorization, Content-Type, and Accept: text/event-stream

### Requirement: Stream options support
The system SHALL support stream_options parameter with include_usage field.

#### Scenario: Include usage option
- **WHEN** stream_options.include_usage is set to true
- **THEN** request includes stream_options in body

#### Scenario: Usage in final chunk
- **WHEN** include_usage is true
- **THEN** final stream chunk includes Usage field

### Requirement: Streaming response handling
The system SHALL parse and deliver streaming chunks as they arrive.

#### Scenario: First chunk with role
- **WHEN** first chunk arrives
- **THEN** Delta includes role field (assistant)

#### Scenario: Content chunks
- **WHEN** subsequent chunks arrive
- **THEN** Delta includes content field with token text

#### Scenario: Finish reason chunk
- **WHEN** stream completes
- **THEN** final chunk includes finish_reason in choices

### Requirement: Error handling for streaming
The system SHALL handle API errors specific to streaming requests.

#### Scenario: Authentication error
- **WHEN** API returns 401 Unauthorized
- **THEN** error is returned before stream starts

#### Scenario: Bad request error
- **WHEN** API returns 400 Bad Request
- **THEN** error is returned before stream starts

#### Scenario: Server error during stream
- **WHEN** API returns 5xx during streaming
- **THEN** error is sent on error channel

### Requirement: Integration with types
The system SHALL use types from types-definitions change for streaming.

#### Scenario: Request type compatibility
- **WHEN** CreateStream is called
- **THEN** ChatCompletionRequest type is used (same as non-streaming)

#### Scenario: Response type compatibility
- **WHEN** stream chunk is received
- **THEN** StreamResponse type is compatible with ChatCompletion type
