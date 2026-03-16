## ADDED Requirements

### Requirement: Core message types
The system SHALL define types for chat messages including role enumeration and message structure with content and optional fields.

#### Scenario: System message creation
- **WHEN** a developer creates a system message
- **THEN** the Message struct accepts MessageRoleSystem and string content

#### Scenario: User message creation
- **WHEN** a developer creates a user message
- **THEN** the Message struct accepts MessageRoleUser and string content

#### Scenario: Assistant message creation
- **WHEN** a developer creates an assistant message
- **THEN** the Message struct accepts MessageRoleAssistant with content and optional fields

### Requirement: Chat completion request types
The system SHALL define ChatCompletionRequest type with all parameters supported by Cerebras Cloud API including model, messages, temperature, max_tokens, and optional parameters.

#### Scenario: Minimal request
- **WHEN** a developer creates a minimal chat completion request
- **THEN** only model and messages fields are required

#### Scenario: Full request with optional parameters
- **WHEN** a developer creates a request with optional parameters
- **THEN** temperature, max_tokens, top_p, stop, n, and other optional fields can be set via pointers

#### Scenario: Request JSON serialization
- **WHEN** a ChatCompletionRequest is marshaled to JSON
- **THEN** the output matches Cerebras Cloud API expected format with snake_case field names

### Requirement: Chat completion response types
The system SHALL define ChatCompletion response type with choices, usage, time_info, and metadata fields.

#### Scenario: Response with single choice
- **WHEN** API returns a chat completion with one choice
- **THEN** ChatCompletion.Choices contains one element with message and finish_reason

#### Scenario: Response with usage information
- **WHEN** API returns usage metrics
- **THEN** ChatCompletion.Usage contains prompt_tokens, completion_tokens, and total_tokens

#### Scenario: Response with time metrics
- **WHEN** API returns timing information
- **THEN** ChatCompletion.TimeInfo contains created_at, started_at, completed_at, and total_time

### Requirement: Choice and message content types
The system SHALL define Choice type with message, finish_reason, and index fields, plus message content structure.

#### Scenario: Choice with text content
- **WHEN** a choice contains a text response
- **THEN** Message.Content contains the assistant's text response

#### Scenario: Choice with finish reason
- **WHEN** a choice completes
- **THEN** Choice.FinishReason indicates stop, length, tool_calls, or other termination reason

### Requirement: Usage metrics types
The system SHALL define Usage type with token counts and optional detailed breakdowns.

#### Scenario: Basic usage tracking
- **WHEN** API returns token usage
- **THEN** Usage includes prompt_tokens, completion_tokens, and total_tokens

#### Scenario: Detailed usage breakdown
- **WHEN** API returns detailed token breakdown
- **THEN** Usage may include prompt_tokens_details and completion_tokens_details with sub-fields

### Requirement: Model types
The system SHALL define Model and ModelList types for the Models API endpoint.

#### Scenario: Model representation
- **WHEN** a model is retrieved
- **THEN** Model type includes id, object, created, and owned_by fields

#### Scenario: Model list representation
- **WHEN** multiple models are listed
- **THEN** ModelList type includes object, data array, and optional pagination fields

### Requirement: Streaming types
The system SHALL define StreamResponse, Delta, and StreamOptions types for SSE streaming.

#### Scenario: Stream chunk representation
- **WHEN** API returns a streaming chunk
- **THEN** StreamResponse includes id, object, created, model, choices with delta, and optional usage

#### Scenario: Delta content
- **WHEN** a streaming chunk contains new content
- **THEN** Delta includes role (first chunk) and/or content (subsequent chunks)

#### Scenario: Stream options configuration
- **WHEN** a developer configures streaming
- **THEN** StreamOptions can specify include_usage to receive usage in final chunk

### Requirement: Tool calling types
The system SHALL define Tool, ToolChoice, Function, and related types for function calling.

#### Scenario: Tool definition
- **WHEN** a developer defines a tool
- **THEN** Tool type includes type="function" and Function with name, description, and parameters

#### Scenario: Tool choice specification
- **WHEN** a developer specifies tool choice
- **THEN** ToolChoice can be "none", "auto", "required", or specific function name

#### Scenario: Tool call in response
- **WHEN** model invokes a tool
- **THEN** Message includes ToolCalls array with function name and arguments

### Requirement: Response format types
The system SHALL define ResponseFormat type for controlling output format (text, json_object, json_schema).

#### Scenario: JSON object format
- **WHEN** a developer requests JSON output
- **THEN** ResponseFormat.Type is set to "json_object"

#### Scenario: JSON schema format
- **WHEN** a developer requests structured JSON
- **THEN** ResponseFormat includes json_schema with name, description, and schema definition

### Requirement: Common parameter types
The system SHALL define common parameter types used across multiple endpoints including frequency_penalty, presence_penalty, logprobs, logit_bias, seed, user, and service_tier.

#### Scenario: Penalty parameters
- **WHEN** a developer sets frequency or presence penalty
- **THEN** parameters accept float64 values between -2.0 and 2.0

#### Scenario: Logprobs configuration
- **WHEN** a developer requests log probabilities
- **THEN** LogprobsRequest includes top_logprobs count and other options

### Requirement: Error response types
The system SHALL define APIError response type matching the error hierarchy in the errors package.

#### Scenario: Error response structure
- **WHEN** API returns an error
- **THEN** ErrorResponse includes error object with message, type, param, and code fields
