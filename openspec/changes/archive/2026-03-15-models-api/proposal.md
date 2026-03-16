## Why

The Go SDK lacks the Models API endpoint implementation, which is essential for discovering available models and retrieving model metadata. This change adds the Models API enabling users to list available models and retrieve specific model details.

## What Changes

- Implement `ModelsService` with `List` and `Retrieve` methods
- Add `GET /v1/models` endpoint support for listing models
- Add `GET /v1/models/{model_id}` endpoint support for retrieving a specific model
- Return typed `ModelList` and `Model` responses
- Integrate with existing HTTP transport and error handling
- Add comprehensive examples and documentation

## Capabilities

### New Capabilities
- `models-list`: List available models via GET /v1/models endpoint
- `models-retrieve`: Retrieve specific model details via GET /v1/models/{model_id} endpoint

### Modified Capabilities
- None

## Impact

- **Affected code**: New service in `pkg/cerebras/models.go`
- **Dependencies**: Requires types from `types-definitions` change
- **APIs**: New public API `client.Models.List(ctx)` and `client.Models.Retrieve(ctx, modelID)`
- **Systems**: Enables model discovery and validation for chat completions
