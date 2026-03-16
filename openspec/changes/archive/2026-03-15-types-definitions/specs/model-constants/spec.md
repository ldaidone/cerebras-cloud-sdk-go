## ADDED Requirements

### Requirement: Model identifier constants
The system SHALL define string constants for all supported model identifiers to prevent typos and provide IDE autocomplete.

#### Scenario: Llama 3.1 8B model constant
- **WHEN** a developer references the Llama 3.1 8B model
- **THEN** the constant Llama31_8b or ModelLlama31_8b provides the value "llama3.1-8b"

#### Scenario: Llama 3.1 70B model constant
- **WHEN** a developer references the Llama 3.1 70B model
- **THEN** the constant Llama31_70b or ModelLlama31_70b provides the value "llama3.1-70b"

#### Scenario: GPT-oss 120B model constant
- **WHEN** a developer references the GPT-oss 120B model
- **THEN** the constant GptOss120b or ModelGptOss120b provides the value "gpt-oss-120b"

#### Scenario: Qwen 3 235B model constant
- **WHEN** a developer references the Qwen 3 235B model
- **THEN** the constant Qwen3_235b or ModelQwen3_235b provides the value "qwen-3-235b-a22b"

#### Scenario: ZAI GLM 4.7 model constant
- **WHEN** a developer references the ZAI GLM 4.7 model
- **THEN** the constant ZaiGlm47 or ModelZaiGlm47 provides the value "zai-glm-4.7"

### Requirement: Model constant organization
The system SHALL organize model constants in a clear, discoverable manner within the package.

#### Scenario: Constant naming convention
- **WHEN** a developer searches for model constants
- **THEN** constants follow consistent naming (PascalCase, descriptive names matching model identifiers)

#### Scenario: Constant documentation
- **WHEN** a developer views constant documentation
- **THEN** each constant includes GoDoc comment describing the model
