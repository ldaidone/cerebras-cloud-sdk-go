// Package cerebras provides the Models API for the Cerebras Cloud SDK.
//
// The Models API allows you to list and retrieve information about
// available models on the Cerebras Cloud platform.
//
// Basic usage:
//
//	client := cerebras.NewClient(cerebras.WithAPIKey("sk-..."))
//
//	// List all models
//	models, err := client.Models.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, model := range models.Data {
//		fmt.Printf("Model: %s (owned by %s)\n", model.ID, model.OwnedBy)
//	}
//
//	// Retrieve a specific model
//	model, err := client.Models.Retrieve(ctx, cerebras.Llama31_8b)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Model: %s, Created: %d\n", model.ID, model.Created)
package cerebras

import (
	"context"
	"fmt"
	"net/http"

	cerebraserrors "github/ldaidone/cerebras-cloud-sdk-go/internal/errors"
	"github/ldaidone/cerebras-cloud-sdk-go/internal/transport"
)

// ModelsService provides access to the Models API.
//
// This service allows you to list available models and retrieve details
// about specific models.
type ModelsService struct {
	client    *Client
	transport *transport.Transport
}

// List retrieves a list of all available models.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//
// Returns:
//   - *ModelList: The list of models
//   - error: Any error that occurred during the request
//
// Example:
//
//	models, err := client.Models.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, model := range models.Data {
//		fmt.Printf("Model: %s (owned by %s)\n", model.ID, model.OwnedBy)
//	}
func (s *ModelsService) List(ctx context.Context) (*ModelList, error) {
	// Make the HTTP GET request
	var result ModelList
	err := s.client.Get(ctx, "/v1/models", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Retrieve retrieves details about a specific model.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - modelID: The model identifier (e.g., "llama3.1-8b", Llama31_8b)
//
// Returns:
//   - *Model: The model details
//   - error: Any error that occurred during the request
//
// Example:
//
//	model, err := client.Models.Retrieve(ctx, cerebras.Llama31_8b)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Model: %s, Created: %d, Owned by: %s\n",
//		model.ID, model.Created, model.OwnedBy)
func (s *ModelsService) Retrieve(ctx context.Context, modelID string) (*Model, error) {
	// Validate modelID
	if modelID == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "modelID parameter is required"))
	}

	// Make the HTTP GET request
	var result Model
	path := fmt.Sprintf("/v1/models/%s", modelID)
	err := s.client.Get(ctx, path, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
