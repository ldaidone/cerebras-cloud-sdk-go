## ADDED Requirements

### Requirement: Response format types
The system SHALL define types for response format specification including ResponseFormat and ResponseFormatType.

#### Scenario: ResponseFormatType enum
- **WHEN** a developer specifies response format
- **THEN** ResponseFormatType accepts "text", "json_object", or "json_schema" values

#### Scenario: ResponseFormat struct definition
- **WHEN** a developer defines response format
- **THEN** ResponseFormat struct includes Type and optional JSONSchema fields

#### Scenario: ResponseFormat constants
- **WHEN** a developer uses response format constants
- **THEN** ResponseFormatText, ResponseFormatJSONObject, ResponseFormatJSONSchema are available

#### Scenario: ResponseFormat JSON serialization
- **WHEN** ResponseFormat is marshaled to JSON
- **THEN** output matches Cerebras Cloud API expected format

### Requirement: JSON Schema types
The system SHALL define types for JSON schema response format including JSONSchema struct.

#### Scenario: JSONSchema struct definition
- **WHEN** a developer defines JSON schema
- **THEN** JSONSchema struct includes Name, Description, Schema, and optional Strict fields

#### Scenario: Schema definition
- **WHEN** a developer defines JSON schema structure
- **THEN** Schema field accepts JSON Schema as map[string]interface{}

#### Scenario: Strict mode
- **WHEN** a developer enables strict mode
- **THEN** Strict field accepts boolean pointer

#### Scenario: JSONSchema JSON serialization
- **WHEN** JSONSchema is marshaled to JSON
- **THEN** output matches API expected format with required name and schema fields

### Requirement: Helper functions
The system SHALL provide helper functions for common response format patterns.

#### Scenario: JSON object helper
- **WHEN** a developer uses ResponseFormatJSON helper
- **THEN** helper returns ResponseFormat with type="json_object"

#### Scenario: JSON schema helper
- **WHEN** a developer uses ResponseFormatJSONWithSchema helper
- **THEN** helper returns ResponseFormat with type="json_schema" and provided schema
