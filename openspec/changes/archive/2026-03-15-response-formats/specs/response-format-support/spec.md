## ADDED Requirements

### Requirement: Chat completions response format parameter
The system SHALL support response_format parameter in chat completions requests.

#### Scenario: Response format parameter
- **WHEN** a developer includes response_format in request
- **THEN** ChatCompletionRequest.ResponseFormat accepts *ResponseFormat

#### Scenario: JSON object mode
- **WHEN** response_format type is "json_object"
- **THEN** model returns valid JSON object in response content

#### Scenario: JSON schema mode
- **WHEN** response_format type is "json_schema" with schema
- **THEN** model returns JSON matching the provided schema

#### Scenario: Request JSON serialization
- **WHEN** request with response_format is marshaled
- **THEN** JSON includes response_format field with proper structure

### Requirement: Response format behavior
The system SHALL handle response format according to API behavior.

#### Scenario: JSON object response
- **WHEN** json_object mode is requested
- **THEN** response content contains valid JSON string

#### Scenario: JSON schema response
- **WHEN** json_schema mode is requested
- **THEN** response content matches provided schema structure

#### Scenario: Text mode (default)
- **WHEN** no response_format is specified
- **THEN** model returns natural language text

### Requirement: Error handling
The system SHALL handle errors specific to response format.

#### Scenario: Invalid JSON schema
- **WHEN** API returns 400 for invalid schema
- **THEN** BadRequestError is returned with error details

#### Scenario: Schema validation error
- **WHEN** API returns error for schema issues
- **THEN** appropriate error is returned with details

### Requirement: Integration with other features
The system SHALL support response format with other chat completion features.

#### Scenario: Response format with tools
- **WHEN** both response_format and tools are specified
- **THEN** request is sent with both parameters

#### Scenario: Response format with streaming
- **WHEN** response_format is used with streaming
- **THEN** streaming chunks contain formatted content
