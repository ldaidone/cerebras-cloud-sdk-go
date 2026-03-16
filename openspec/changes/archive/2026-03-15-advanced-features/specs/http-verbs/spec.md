## ADDED Requirements

### Requirement: HTTP verb helpers
The system SHALL provide HTTP verb helper methods on the client for custom API calls.

#### Scenario: Get method
- **WHEN** a developer calls client.Get
- **THEN** HTTP GET request is sent to specified path with response deserialized

#### Scenario: Post method
- **WHEN** a developer calls client.Post
- **THEN** HTTP POST request is sent with body serialized to JSON

#### Scenario: Put method
- **WHEN** a developer calls client.Put
- **THEN** HTTP PUT request is sent with body serialized to JSON

#### Scenario: Delete method
- **WHEN** a developer calls client.Delete
- **THEN** HTTP DELETE request is sent with response deserialized

### Requirement: HTTP verb integration
The system SHALL integrate HTTP verb helpers with existing transport layer.

#### Scenario: Authentication
- **WHEN** HTTP verb methods are called
- **THEN** requests include Authorization header with API key

#### Scenario: Error handling
- **WHEN** HTTP verb request fails
- **THEN** typed errors from error hierarchy are returned

#### Scenario: Context propagation
- **WHEN** context is passed to HTTP verb methods
- **THEN** context is propagated for cancellation and timeout

### Requirement: WithOptions method
The system SHALL provide WithOptions method for client configuration chaining.

#### Scenario: Client cloning
- **WHEN** a developer calls WithOptions
- **THEN** new client instance is returned with modified options

#### Scenario: Option inheritance
- **WHEN** new client is created with WithOptions
- **THEN** original client settings are inherited unless overridden

#### Scenario: Non-mutating operation
- **WHEN** WithOptions is called
- **THEN** original client remains unchanged
