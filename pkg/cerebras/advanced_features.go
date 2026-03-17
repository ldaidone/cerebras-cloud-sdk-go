// Package cerebras provides advanced features for the Cerebras Cloud SDK.
package cerebras

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cerebraserrors "github.com/ldaidone/cerebras-cloud-sdk-go/internal/errors"
)

// ServiceTier represents the service tier for request prioritization.
type ServiceTier string

const (
	// ServiceTierAuto automatically selects the service tier.
	ServiceTierAuto ServiceTier = "auto"
	// ServiceTierDefault uses the default service tier.
	ServiceTierDefault ServiceTier = "default"
	// ServiceTierFlex uses the flex service tier (lower cost, variable latency).
	ServiceTierFlex ServiceTier = "flex"
	// ServiceTierPriority uses the priority service tier (higher cost, lower latency).
	ServiceTierPriority ServiceTier = "priority"
)

// ReasoningEffort represents the reasoning effort for reasoning models.
type ReasoningEffort string

const (
	// ReasoningEffortLow uses minimal reasoning.
	ReasoningEffortLow ReasoningEffort = "low"
	// ReasoningEffortMedium uses moderate reasoning.
	ReasoningEffortMedium ReasoningEffort = "medium"
	// ReasoningEffortHigh uses extensive reasoning.
	ReasoningEffortHigh ReasoningEffort = "high"
)

// ReasoningFormat represents how reasoning is returned.
type ReasoningFormat string

const (
	// ReasoningFormatNone uses the model default.
	ReasoningFormatNone ReasoningFormat = "none"
	// ReasoningFormatParsed returns reasoning in the reasoning field.
	ReasoningFormatParsed ReasoningFormat = "parsed"
	// ReasoningFormatTextParsed returns reasoning like parsed but logprobs not separated.
	ReasoningFormatTextParsed ReasoningFormat = "text_parsed"
	// ReasoningFormatRaw returns reasoning in content with special tokens.
	ReasoningFormatRaw ReasoningFormat = "raw"
	// ReasoningFormatHidden does not return reasoning.
	ReasoningFormatHidden ReasoningFormat = "hidden"
)

// APIResponse wraps a response with access to raw HTTP response data.
// This is useful for debugging and accessing headers.
type APIResponse struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Headers contains the response headers.
	Headers map[string][]string
	// Body contains the raw response body.
	Body []byte
}

// WithRawResponse executes a request and returns the raw response.
// This is useful for debugging or accessing response headers.
func (s *ChatCompletionsService) WithRawResponse(
	ctx context.Context,
	model string,
	messages []Message,
	opts ...ChatCompletionOption,
) (*APIResponse, error) {
	// Build request
	req := &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	for _, opt := range opts {
		opt(req)
	}

	// Marshal request body
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to marshal request body: %w", err))
	}

	// Create HTTP request
	reqURL := s.transport.GetBaseURL() + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to create request: %w", err))
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)

	// Execute request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to execute request: %w", err))
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to read response body: %w", err))
	}

	// Build APIResponse
	apiResp := &APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
	}

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return apiResp, cerebraserrors.MapStatusCode(resp.StatusCode, string(body))
	}

	return apiResp, nil
}
