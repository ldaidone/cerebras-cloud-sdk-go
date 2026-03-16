## ADDED Requirements

### Requirement: Tool definition types
The system SHALL define types for tool definitions including Tool, Function, and parameter schemas.

#### Scenario: Tool struct definition
- **WHEN** a developer defines a tool
- **THEN** Tool struct includes Type (string) and Function (Function struct) fields

#### Scenario: Function struct definition
- **WHEN** a developer defines a function
- **THEN** Function struct includes Name, Description, and Parameters fields

#### Scenario: JSON Schema parameters
- **WHEN** a developer defines function parameters
- **THEN** Parameters field accepts JSON Schema as map[string]interface{}

#### Scenario: Tool JSON serialization
- **WHEN** Tool is marshaled to JSON
- **THEN** output matches Cerebras Cloud API expected format with type="function"

### Requirement: Tool choice types
The system SHALL define types for tool choice specification including string and object forms.

#### Scenario: String tool choice
- **WHEN** a developer specifies tool choice as string
- **THEN** "none", "auto", or "required" values are accepted

#### Scenario: Object tool choice
- **WHEN** a developer specifies tool choice as object
- **THEN** ToolChoiceFunction struct with type and function fields is supported

#### Scenario: Tool choice function name
- **WHEN** a developer specifies a specific function
- **THEN** ToolChoiceFunctionName struct with name field is supported

#### Scenario: Tool choice JSON serialization
- **WHEN** ToolChoice is marshaled to JSON
- **THEN** output matches API expected format (string or object)

### Requirement: Tool call response types
The system SHALL define types for tool calls in model responses including ToolCall and FunctionCall.

#### Scenario: ToolCall struct definition
- **WHEN** model returns tool calls
- **THEN** ToolCall struct includes ID, Type, and Function fields

#### Scenario: FunctionCall struct definition
- **WHEN** model invokes a function
- **THEN** FunctionCall struct includes Name and Arguments fields

#### Scenario: Function arguments as JSON string
- **WHEN** model returns function arguments
- **THEN** Arguments field contains JSON string (not parsed object)

#### Scenario: ToolCall in message
- **WHEN** message includes tool calls
- **THEN** Message.ToolCalls array contains ToolCall structs

### Requirement: Helper functions
The system SHALL provide helper functions for common tool definition patterns.

#### Scenario: DefineFunction helper
- **WHEN** a developer uses DefineFunction helper
- **THEN** helper creates Tool struct with proper structure

#### Scenario: JSON Schema helpers
- **WHEN** a developer defines JSON Schema
- **THEN** helpers assist with property and required field definition
