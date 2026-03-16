## ADDED Requirements

### Requirement: Chat completions service
The system SHALL provide a ChatCompletionsService accessible via client.ChatCompletions with a Create method for non-streaming chat completions.

#### Scenario: Service access
- **WHEN** a user initializes the client
- **THEN** client.ChatCompletions provides access to the chat completions service

#### Scenario: Create method signature
- **WHEN** a user calls the Create method
- **THEN** the method accepts context, model, messages, and optional parameters

### Requirement: POST /v1/chat/completions endpoint
The system SHALL implement HTTP POST to /v1/chat/completions with proper request/response handling.

#### Scenario: Request endpoint
- **WHEN** Create method is called
- **THEN** HTTP POST request is sent to /v1/chat/completions

#### Scenario: Request headers
- **WHEN** request is sent
- **THEN** headers include Authorization (Bearer token), Content-Type (application/json), and User-Agent

#### Scenario: Request body serialization
- **WHEN** request is sent
- **THEN** ChatCompletionRequest is serialized to JSON with snake_case field names

### Requirement: Required parameters
The system SHALL require model and messages parameters for chat completions requests.

#### Scenario: Model parameter
- **WHEN** a request is made without a model
- **THEN** the SDK returns an error indicating model is required

#### Scenario: Messages parameter
- **WHEN** a request is made without messages
- **THEN** the SDK returns an error indicating messages are required

#### Scenario: Empty messages array
- **WHEN** a request is made with an empty messages array
- **THEN** the SDK returns an error indicating at least one message is required

### Requirement: Optional parameters support
The system SHALL support all optional parameters defined in the Cerebras Cloud API including temperature, max_tokens, top_p, stop, n, frequency_penalty, presence_penalty, logprobs, logit_bias, seed, user, and service_tier.

#### Scenario: Temperature parameter
- **WHEN** temperature is set
- **THEN** the value is included in the request body

#### Scenario: Max tokens parameter
- **WHEN** max_tokens is set
- **THEN** the value is included in the request body

#### Scenario: Top P parameter
- **WHEN** top_p is set
- **THEN** the value is included in the request body

#### Scenario: Stop sequences parameter
- **WHEN** stop sequences are set
- **THEN** the values are included in the request body as an array or string

#### Scenario: Number of choices parameter
- **WHEN** n (number of choices) is set
- **THEN** the value is included in the request body

### Requirement: Response handling
The system SHALL parse and return ChatCompletion response with choices, usage, and metadata.

#### Scenario: Successful response
- **WHEN** API returns 200 OK
- **THEN** ChatCompletion struct is returned with all fields populated

#### Scenario: Response with single choice
- **WHEN** API returns one choice
- **THEN** Choices array contains one element with message and finish_reason

#### Scenario: Response with multiple choices
- **WHEN** API returns n>1 choices
- **THEN** Choices array contains n elements, each with unique index

#### Scenario: Response with usage information
- **WHEN** API returns usage metrics
- **THEN** Usage field contains prompt_tokens, completion_tokens, and total_tokens

#### Scenario: Response with time metrics
- **WHEN** API returns timing information
- **THEN** TimeInfo field contains created_at, started_at, completed_at, and total_time

### Requirement: Error handling
The system SHALL handle API errors and return typed errors from the error hierarchy.

#### Scenario: Authentication error
- **WHEN** API returns 401 Unauthorized
- **THEN** AuthenticationError is returned

#### Scenario: Bad request error
- **WHEN** API returns 400 Bad Request
- **THEN** BadRequestError is returned with error details

#### Scenario: Rate limit error
- **WHEN** API returns 429 Too Many Requests
- **THEN** RateLimitError is returned with retry information

#### Scenario: Server error
- **WHEN** API returns 5xx error
- **THEN** InternalServerError or specific 5xx error is returned

### Requirement: Context propagation
The system SHALL propagate context for cancellation and timeout handling.

#### Scenario: Context cancellation
- **WHEN** context is cancelled during request
- **THEN** request is aborted and context.Canceled error is returned

#### Scenario: Context timeout
- **WHEN** context deadline exceeds during request
- **THEN** request is aborted and context.DeadlineExceeded error is returned

### Requirement: Retry behavior
The system SHALL apply automatic retry logic for retryable errors.

#### Scenario: Retry on 5xx error
- **WHEN** API returns 503 Service Unavailable
- **THEN** request is retried according to retry policy

#### Scenario: Retry on rate limit
- **WHEN** API returns 429 with Retry-After header
- **THEN** request is retried after the specified delay

#### Scenario: No retry on client error
- **WHEN** API returns 400 Bad Request
- **THEN** request is not retried and error is returned immediately
