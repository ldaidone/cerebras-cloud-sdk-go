## ADDED Requirements

### Requirement: Text completions streaming service
The system SHALL provide streaming support for text completions via CreateStream method.

#### Scenario: CreateStream method signature
- **WHEN** a user calls CreateStream method
- **THEN** the method accepts context and TextCompletionRequest, returns stream channel and error

#### Scenario: Stream parameter
- **WHEN** request is sent
- **THEN** request body includes stream=true

### Requirement: SSE streaming for text completions
The system SHALL use SSE streaming infrastructure for text completions streaming.

#### Scenario: Request headers
- **WHEN** streaming request is sent
- **THEN** headers include Accept: text/event-stream

#### Scenario: SSE event parsing
- **WHEN** API returns SSE events
- **THEN** events are parsed into TextCompletionStreamResponse chunks

### Requirement: Stream response structure
The system SHALL define TextCompletionStreamResponse for streaming chunks.

#### Scenario: Stream chunk structure
- **WHEN** API returns streaming chunk
- **THEN** chunk includes id, object, created, model, and choices

#### Scenario: Text delta
- **WHEN** chunk contains new text
- **THEN** choices include text delta

#### Scenario: Final chunk with usage
- **WHEN** include_usage is true
- **THEN** final chunk includes Usage field

### Requirement: Error handling for streaming
The system SHALL handle errors during text completions streaming.

#### Scenario: Stream error channel
- **WHEN** streaming request is initiated
- **THEN** error channel is provided for stream errors

#### Scenario: HTTP error during stream
- **WHEN** HTTP connection fails during streaming
- **THEN** error is sent on error channel and stream channel is closed

### Requirement: Context and cleanup
The system SHALL support context cancellation and proper cleanup for streaming.

#### Scenario: Context cancellation
- **WHEN** context is cancelled during streaming
- **THEN** HTTP request is aborted and channel is closed

#### Scenario: Response body closing
- **WHEN** stream completes or errors
- **THEN** HTTP response body is closed
