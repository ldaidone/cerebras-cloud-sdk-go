## ADDED Requirements

### Requirement: Reasoning parameter types
The system SHALL define types for reasoning model parameters including reasoning_effort, reasoning_format, clear_thinking, and disable_reasoning.

#### Scenario: Reasoning effort constants
- **WHEN** a developer specifies reasoning effort
- **THEN** constants ReasoningEffortLow, ReasoningEffortMedium, ReasoningEffortHigh are available

#### Scenario: Reasoning format constants
- **WHEN** a developer specifies reasoning format
- **THEN** constants for reasoning format values are available

#### Scenario: Reasoning parameter types
- **WHEN** a developer uses reasoning parameters
- **THEN** parameters accept pointer types for optional specification

### Requirement: Reasoning parameter support
The system SHALL support reasoning parameters in chat completions requests.

#### Scenario: Reasoning effort parameter
- **WHEN** a developer includes reasoning_effort in request
- **THEN** ChatCompletionRequest.ReasoningEffort accepts *string

#### Scenario: Reasoning format parameter
- **WHEN** a developer includes reasoning_format in request
- **THEN** ChatCompletionRequest.ReasoningFormat accepts *string

#### Scenario: Clear thinking parameter
- **WHEN** a developer includes clear_thinking in request
- **THEN** ChatCompletionRequest.ClearThinking accepts *bool

#### Scenario: Disable reasoning parameter
- **WHEN** a developer includes disable_reasoning in request
- **THEN** ChatCompletionRequest.DisableReasoning accepts *bool

### Requirement: Reasoning response handling
The system SHALL handle reasoning content in responses.

#### Scenario: Reasoning content in response
- **WHEN** model returns reasoning content
- **THEN** response includes reasoning field in message

#### Scenario: Reasoning tokens in usage
- **WHEN** model uses reasoning tokens
- **THEN** usage includes reasoning_tokens count

### Requirement: Reasoning model compatibility
The system SHALL document reasoning parameter compatibility.

#### Scenario: Reasoning-capable models
- **WHEN** reasoning parameters are used with compatible models
- **THEN** parameters are applied correctly

#### Scenario: Non-reasoning models
- **WHEN** reasoning parameters are used with incompatible models
- **THEN** API returns appropriate error or ignores parameters
