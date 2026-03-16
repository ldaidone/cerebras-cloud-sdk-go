## ADDED Requirements

### Requirement: API response wrapper type
The system SHALL define APIResponse wrapper type for raw response access.

#### Scenario: APIResponse struct definition
- **WHEN** a developer uses raw response access
- **THEN** APIResponse type includes Data, StatusCode, Headers, and RawBody fields

#### Scenario: Generic type parameter
- **WHEN** a developer specifies response type
- **THEN** APIResponse accepts generic type parameter for Data field

#### Scenario: Response metadata access
- **WHEN** a developer accesses response metadata
- **THEN** StatusCode, Headers, and RawBody are accessible

### Requirement: Raw response methods
The system SHALL provide methods for accessing raw HTTP response.

#### Scenario: CreateWithResponse method
- **WHEN** a developer calls CreateWithResponse
- **THEN** method returns *APIResponse[ChatCompletion] with full response data

#### Scenario: Raw body access
- **WHEN** a developer accesses RawBody
- **THEN** raw JSON response bytes are available

#### Scenario: Header access
- **WHEN** a developer accesses Headers
- **THEN** HTTP response headers are available for inspection

### Requirement: Debugging use cases
The system SHALL support debugging and inspection scenarios.

#### Scenario: Status code inspection
- **WHEN** a developer needs to check status code
- **THEN** StatusCode field provides HTTP status code

#### Scenario: Rate limit header access
- **WHEN** a developer checks rate limit headers
- **THEN** Headers include X-RateLimit-Limit and X-RateLimit-Remaining

#### Scenario: Response logging
- **WHEN** a developer logs raw response
- **THEN** RawBody can be logged for debugging
