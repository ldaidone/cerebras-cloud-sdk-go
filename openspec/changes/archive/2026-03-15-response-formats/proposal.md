## Why

The Go SDK lacks response format controls, which enable users to request structured JSON output from models. This change adds response_format parameter support allowing users to specify JSON object or JSON schema output for reliable structured data extraction.

## What Changes

- Add ResponseFormat type with type field (text, json_object, json_schema)
- Add JSONSchema type for schema-constrained output
- Integrate response_format parameter with chat completions
- Add validation helpers for JSON schema response format
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `response-format-types`: Type definitions for response_format including ResponseFormat and JSONSchema
- `response-format-support`: Response format integration with chat completions endpoint

### Modified Capabilities
- None

## Impact

- **Affected code**: New types in `pkg/cerebras/`, chat completions integration
- **Dependencies**: Requires types from `types-definitions` change
- **APIs**: Extended ChatCompletionRequest with response_format field
- **Systems**: Enables structured JSON output for reliable parsing
