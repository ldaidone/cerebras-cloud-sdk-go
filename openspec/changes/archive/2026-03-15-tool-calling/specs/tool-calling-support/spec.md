## ADDED Requirements

### Requirement: Chat completions tool parameters
The system SHALL support tools, tool_choice, and parallel_tool_calls parameters in chat completions requests.

#### Scenario: Tools parameter
- **WHEN** a developer includes tools in request
- **THEN** ChatCompletionRequest.Tools accepts array of Tool structs

#### Scenario: Tool choice parameter
- **WHEN** a developer specifies tool choice
- **THEN** ChatCompletionRequest.ToolChoice accepts string or ToolChoice object

#### Scenario: Parallel tool calls parameter
- **WHEN** a developer sets parallel_tool_calls
- **THEN** ChatCompletionRequest.ParallelToolCalls accepts boolean pointer

#### Scenario: Request JSON serialization
- **WHEN** request with tools is marshaled
- **THEN** JSON includes tools, tool_choice, and parallel_tool_calls fields

### Requirement: Tool calling in responses
The system SHALL parse and return tool calls in chat completion responses.

#### Scenario: Response with tool calls
- **WHEN** API returns response with tool calls
- **THEN** ChatCompletion.Choices[0].Message.ToolCalls is populated

#### Scenario: Multiple tool calls
- **WHEN** model makes multiple tool calls
- **THEN** ToolCalls array contains all tool calls with unique IDs

#### Scenario: Tool call with arguments
- **WHEN** model invokes function with arguments
- **THEN** FunctionCall.Arguments contains valid JSON string

### Requirement: Tool calling workflow
The system SHALL support the complete tool calling workflow pattern.

#### Scenario: Initial request with tools
- **WHEN** user sends request with tools defined
- **THEN** model may return response with tool calls

#### Scenario: Function execution
- **WHEN** user receives tool calls
- **THEN** user can parse arguments and execute functions

#### Scenario: Submitting function results
- **WHEN** user submits function results
- **THEN** new message with role="tool" and tool_call_id is added to conversation

### Requirement: Integration with streaming
The system SHALL support tool calling in streaming responses.

#### Scenario: Streaming tool calls
- **WHEN** streaming request includes tools
- **THEN** stream chunks may include tool call deltas

#### Scenario: Tool call delta accumulation
- **WHEN** stream returns tool call chunks
- **THEN** user can accumulate deltas to build complete tool call

### Requirement: Error handling
The system SHALL handle errors specific to tool calling.

#### Scenario: Invalid tool definition
- **WHEN** API returns 400 for invalid tool schema
- **THEN** BadRequestError is returned with error details

#### Scenario: Unsupported tool type
- **WHEN** API returns error for unsupported tool type
- **THEN** appropriate error is returned
