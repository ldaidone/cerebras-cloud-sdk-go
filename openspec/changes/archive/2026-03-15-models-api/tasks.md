## 1. Service Structure

- [x] 1.1 Create ModelsService struct in pkg/cerebras/models.go
- [x] 1.2 Add models field to Client struct
- [x] 1.3 Initialize ModelsService in NewClient or constructor

## 2. List Method Implementation

- [x] 2.1 Implement List method signature with context
- [x] 2.2 Implement HTTP GET /v1/models request
- [x] 2.3 Set proper request headers (Authorization, User-Agent)
- [x] 2.4 Handle response deserialization to ModelList struct

## 3. Retrieve Method Implementation

- [x] 3.1 Implement Retrieve method signature with context and modelID
- [x] 3.2 Implement HTTP GET /v1/models/{model_id} request
- [x] 3.3 Handle URL path construction with modelID
- [x] 3.4 Handle response deserialization to Model struct

## 4. Error Handling

- [x] 4.1 Integrate with existing error hierarchy for API errors
- [x] 4.2 Handle 404 Not Found errors for Retrieve
- [x] 4.3 Handle 401 Authentication errors
- [x] 4.4 Handle 5xx Server errors

## 5. Context and Retry Support

- [x] 5.1 Ensure context propagation for cancellation
- [x] 5.2 Verify retry logic applies to models requests
- [x] 5.3 Test context timeout behavior

## 6. Documentation

- [x] 6.1 Add GoDoc comments to ModelsService and methods
- [x] 6.2 Add example usage in doc comments
- [x] 6.3 Create example file showing models list and retrieve usage
