## 1. Tool Type Definitions

- [x] 1.1 Define Tool struct with Type and Function fields
- [x] 1.2 Define Function struct with Name, Description, and Parameters
- [x] 1.3 Define ToolParameters type (map[string]interface{} for JSON Schema)
- [x] 1.4 Add JSON tags to all tool types

## 2. Tool Choice Types

- [x] 2.1 Define ToolChoice interface type (string or object)
- [x] 2.2 Define ToolChoiceFunction struct for object form
- [x] 2.3 Define ToolChoiceFunctionName struct
- [x] 2.4 Implement custom JSON marshaler for ToolChoice if needed

## 3. Tool Call Types

- [x] 3.1 Define ToolCall struct with ID, Type, and Function fields
- [x] 3.2 Define FunctionCall struct with Name and Arguments
- [x] 3.3 Add ToolCalls field to Message struct
- [x] 3.4 Add JSON tags to all tool call types

## 4. Request Integration

- [x] 4.1 Add Tools field to ChatCompletionRequest
- [x] 4.2 Add ToolChoice field to ChatCompletionRequest
- [x] 4.3 Add ParallelToolCalls field to ChatCompletionRequest
- [x] 4.4 Verify JSON serialization with snake_case field names

## 5. Helper Functions

- [x] 5.1 Create DefineFunction helper function
- [x] 5.2 Create JSON Schema helper functions for common patterns
- [x] 5.3 Create WithRequiredParams helper for JSON Schema
- [x] 5.4 Add example tool definitions in documentation

## 6. Documentation

- [x] 6.1 Add GoDoc comments to all tool types
- [x] 6.2 Create example showing tool definition
- [x] 6.3 Create example showing complete tool calling workflow
- [x] 6.4 Document function argument parsing pattern
