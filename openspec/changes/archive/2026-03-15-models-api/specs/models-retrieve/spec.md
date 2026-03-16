## ADDED Requirements

### Requirement: Models retrieve service
The system SHALL provide a ModelsService accessible via client.Models with a Retrieve method for retrieving a specific model by ID.

#### Scenario: Service access
- **WHEN** a user initializes the client
- **THEN** client.Models provides access to the models service

#### Scenario: Retrieve method signature
- **WHEN** a user calls the Retrieve method
- **THEN** the method accepts context and modelID string, returns Model and error

### Requirement: GET /v1/models/{model_id} endpoint
The system SHALL implement HTTP GET to /v1/models/{model_id} with proper request/response handling.

#### Scenario: Request endpoint
- **WHEN** Retrieve method is called with modelID
- **THEN** HTTP GET request is sent to /v1/models/{modelID}

#### Scenario: URL path construction
- **WHEN** request is sent
- **THEN** modelID is properly URL-encoded in the path

#### Scenario: Request headers
- **WHEN** request is sent
- **THEN** headers include Authorization (Bearer token) and User-Agent

#### Scenario: No request body
- **WHEN** request is sent
- **THEN** no request body is included (GET request)

### Requirement: Model retrieval response
The system SHALL parse and return Model response with all model metadata.

#### Scenario: Successful response
- **WHEN** API returns 200 OK
- **THEN** Model struct is returned with all fields populated

#### Scenario: Model structure
- **WHEN** API returns model
- **THEN** Model includes id, object, created, and owned_by fields

#### Scenario: Model ID matches request
- **WHEN** API returns model
- **THEN** Model.ID matches the requested modelID

### Requirement: Not found handling
The system SHALL handle 404 Not Found errors when model doesn't exist.

#### Scenario: Non-existent model
- **WHEN** API returns 404 Not Found
- **THEN** NotFoundError is returned

#### Scenario: Invalid model ID format
- **WHEN** API returns 404 for invalid format
- **THEN** NotFoundError is returned with appropriate message

### Requirement: Error handling
The system SHALL handle API errors and return typed errors from the error hierarchy.

#### Scenario: Authentication error
- **WHEN** API returns 401 Unauthorized
- **THEN** AuthenticationError is returned

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

#### Scenario: No retry on 404
- **WHEN** API returns 404 Not Found
- **THEN** request is not retried and NotFoundError is returned immediately
