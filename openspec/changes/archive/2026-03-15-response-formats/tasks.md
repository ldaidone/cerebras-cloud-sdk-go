## 1. Response Format Types

- [x] 1.1 Define ResponseFormatType string enum with constants
- [x] 1.2 Define ResponseFormat struct with Type and JSONSchema fields
- [x] 1.3 Define JSONSchema struct with Name, Description, Schema, Strict
- [x] 1.4 Add JSON tags to all response format types

## 2. Helper Functions

- [x] 2.1 Create ResponseFormatJSON helper function
- [x] 2.2 Create ResponseFormatJSONWithSchema helper function
- [x] 2.3 Add example JSON Schema definitions in documentation
- [x] 2.4 Create ResponseFormatPlainText helper function

## 3. Request Integration

- [x] 3.1 Add ResponseFormat field to ChatCompletionRequest
- [x] 3.2 Verify JSON serialization with snake_case field names
- [x] 3.3 Test json_object mode request/response
- [x] 3.4 Test json_schema mode request/response

## 4. Documentation

- [x] 4.1 Add GoDoc comments to response format types
- [x] 4.2 Create example showing JSON object response format
- [x] 4.3 Create example showing JSON schema response format
- [x] 4.4 Document JSON Schema structure expectations
