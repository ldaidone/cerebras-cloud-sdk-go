## ADDED Requirements

### Requirement: Service tier types
The system SHALL define types for service tier specification including ServiceTier enum.

#### Scenario: ServiceTier enum definition
- **WHEN** a developer specifies service tier
- **THEN** ServiceTier accepts "auto", "flex", or "priority" values

#### Scenario: Service tier constants
- **WHEN** a developer uses service tier constants
- **THEN** ServiceTierAuto, ServiceTierFlex, ServiceTierPriority are available

#### Scenario: Service tier JSON serialization
- **WHEN** ServiceTier is marshaled to JSON
- **THEN** output matches Cerebras Cloud API expected format

### Requirement: Service tier parameter support
The system SHALL support service_tier parameter in chat completions requests.

#### Scenario: Service tier parameter
- **WHEN** a developer includes service_tier in request
- **THEN** ChatCompletionRequest.ServiceTier accepts *ServiceTier

#### Scenario: Auto tier (default)
- **WHEN** service_tier is not specified
- **THEN** API uses default auto tier behavior

#### Scenario: Flex tier
- **WHEN** service_tier is set to "flex"
- **THEN** request uses flex tier processing

#### Scenario: Priority tier
- **WHEN** service_tier is set to "priority"
- **THEN** request uses priority tier processing

### Requirement: Service tier behavior
The system SHALL handle service tier according to API behavior.

#### Scenario: Tier availability
- **WHEN** requested tier is not available
- **THEN** API returns appropriate error

#### Scenario: Tier pricing
- **WHEN** different tiers have different pricing
- **THEN** usage response may include tier information
