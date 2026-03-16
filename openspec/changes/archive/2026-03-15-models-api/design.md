## Context

The Go SDK has completed its foundation layer and type definitions. The Models API is a standard REST endpoint for discovering available models and retrieving model metadata. This matches the Python SDK's `client.models.list()` and `client.models.retrieve()` pattern.

Key constraints:
- API compatibility with Python SDK method signatures
- RESTful endpoint design (GET /v1/models, GET /v1/models/{model_id})
- Integration with existing transport layer
- Simple request/response pattern (no complex parameters)

## Goals / Non-Goals

**Goals:**
- Implement ModelsService with List and Retrieve methods
- Support GET /v1/models endpoint for listing
- Support GET /v1/models/{model_id} endpoint for retrieval
- Return typed ModelList and Model responses
- Proper error handling with typed errors
- Comprehensive examples and documentation

**Non-Goals:**
- Model creation or deletion (API doesn't support)
- Model filtering or pagination parameters (add if API supports)
- Model metadata updates (API doesn't support)

## Decisions

### 1. Service Structure
**Decision:** Implement as ModelsService struct with method receiver

**Rationale:**
- Matches Python SDK's `client.models` namespace
- Groups related operations logically
- Enables future extension with additional methods

**Structure:**
```go
type ModelsService struct {
    client *Client
}

func (s *ModelsService) List(ctx context.Context) (*ModelList, error)
func (s *ModelsService) Retrieve(ctx context.Context, modelID string) (*Model, error)
```

### 2. Method Naming
**Decision:** Use `List` and `Retrieve` method names

**Rationale:**
- Matches Python SDK naming convention
- Clear and idiomatic Go naming
- Self-documenting method purposes

**Alternatives considered:**
- `Get` instead of `Retrieve`: Less specific, could be confused with GET HTTP method
- `GetModel`: More verbose, inconsistent with Python SDK

### 3. Request/Response Handling
**Decision:** Use existing transport layer's `sendRequest` method

**Rationale:**
- Reuses retry logic, error handling, and JSON marshaling
- Maintains consistency across all endpoints
- Leverages existing context propagation and timeout handling

### 4. URL Path Construction
**Decision:** Use `fmt.Sprintf` for URL path construction

**Rationale:**
- Simple and readable for basic path construction
- No external dependencies needed
- Consistent with Go stdlib patterns

**Implementation:**
```go
// List: GET /v1/models
path := "/v1/models"

// Retrieve: GET /v1/models/{model_id}
path := fmt.Sprintf("/v1/models/%s", modelID)
```

### 5. Error Handling
**Decision:** Return typed errors from existing error hierarchy

**Rationale:**
- Consistent with SDK's error handling approach
- Enables errors.Is/as pattern matching
- Provides clear error messages with HTTP status codes

## Risks / Trade-offs

**[API endpoint changes]** → Cerebras API may change model endpoint structure
- **Mitigation:** Abstract URL construction; easy to update path templates

**[Model ID validation]** → Invalid model IDs may cause confusing errors
- **Mitigation:** Add basic validation (non-empty, valid characters) before request

**[Pagination]** → API may add pagination in future
- **Mitigation:** Design ModelList to support pagination fields (has_more, next_page, etc.)
