## ADDED Requirements

### Requirement: Text completions service
The system SHALL provide a TextCompletionsService accessible via client.TextCompletions with a Create method for non-streaming text completions.

#### Scenario: Service access
- **WHEN** a user initializes the client
- **THEN** client.TextCompletions provides access to the text completions service

#### Scenario: Create method signature
- **WHEN** a user calls the Create method
- **THEN** the method accepts context and TextCompletionRequest, returns TextCompletion and error

### Requirement: POST /v1/completions endpoint
The system SHALL implement HTTP POST to /v1/completions with proper request/response handling.

#### Scenario: Request endpoint
- **WHEN** Create method is called
- **THEN** HTTP POST request is sent to /v1/completions

#### Scenario: Request headers
- **WHEN** request is sent
- **THEN** headers include Authorization (Bearer token), Content-Type (application/json), and User-Agent

#### Scenario: Request body serialization
- **WHEN** request is sent
- **THEN** TextCompletionRequest is serialized to JSON with snake_case field names

### Requirement: Required parameters
The system SHALL require model and prompt parameters for text completions requests.

#### Scenario: Model parameter
- **WHEN** a request is made without a model
- **THEN** the SDK returns an error indicating model is required

#### Scenario: Prompt parameter
- **WHEN** a request is made without a prompt
- **THEN** the SDK returns an error indicating prompt is required

### Requirement: Optional parameters support
The system SHALL support optional parameters including max_tokens, temperature, top_p, stop, n, frequency_penalty, presence_penalty, logprobs, echo, seed, suffix, and max_completion_tokens.

#### Scenario: Temperature parameter
- **WHEN** temperature is set
- **THEN** the value is included in the request body

#### Scenario: Max tokens parameter
- **WHEN** max_tokens is set
- **THEN** the value is included in the request body

#### Scenario: Echo parameter
- **WHEN** echo is set to true
- **THEN** the value is included in the request body

### Requirement: Text completion response
The system SHALL parse and return TextCompletion response with choices, usage, and metadata.

#### Scenario: Successful response
- **WHEN** API returns 200 OK
- **THEN** TextCompletion struct is returned with all fields populated

#### Scenario: Response with single choice
- **WHEN** API returns one choice
- **THEN** Choices array contains one element with text and finish_reason

#### Scenario: Response with multiple choices
- **WHEN** API returns n>1 choices
- **THEN** Choices array contains n elements, each with unique index and text

#### Scenario: Response with usage information
- **WHEN** API returns usage metrics
- **THEN** Usage field contains prompt_tokens, completion_tokens, and total_tokens

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
- **THEN** RateLimitError is returned

### Requirement: Context and retry
The system SHALL propagate context and apply retry logic for text completions.

#### Scenario: Context cancellation
- **WHEN** context is cancelled during request
- **THEN** request is aborted and context.Canceled error is returned

#### Scenario: Retry on 5xx error
- **WHEN** API returns 503 Service Unavailable
- **THEN** request is retried according to retry policy
