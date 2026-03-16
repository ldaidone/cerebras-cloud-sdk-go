## ADDED Requirements

### Requirement: Models list service
The system SHALL provide a ModelsService accessible via client.Models with a List method for retrieving all available models.

#### Scenario: Service access
- **WHEN** a user initializes the client
- **THEN** client.Models provides access to the models service

#### Scenario: List method signature
- **WHEN** a user calls the List method
- **THEN** the method accepts context and returns ModelList and error

### Requirement: GET /v1/models endpoint
The system SHALL implement HTTP GET to /v1/models with proper request/response handling.

#### Scenario: Request endpoint
- **WHEN** List method is called
- **THEN** HTTP GET request is sent to /v1/models

#### Scenario: Request headers
- **WHEN** request is sent
- **THEN** headers include Authorization (Bearer token) and User-Agent

#### Scenario: No request body
- **WHEN** request is sent
- **THEN** no request body is included (GET request)

### Requirement: Model list response
The system SHALL parse and return ModelList response with data array and metadata.

#### Scenario: Successful response
- **WHEN** API returns 200 OK
- **THEN** ModelList struct is returned with all fields populated

#### Scenario: Model list with multiple models
- **WHEN** API returns multiple models
- **THEN** ModelList.Data contains array of Model structs

#### Scenario: Model list structure
- **WHEN** API returns model list
- **THEN** ModelList includes object field (e.g., "list") and data array

#### Scenario: Individual model structure
- **WHEN** API returns model list
- **THEN** each Model includes id, object, created, and owned_by fields

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

#### Scenario: No retry on client error
- **WHEN** API returns 400 Bad Request
- **THEN** request is not retried and error is returned immediately
