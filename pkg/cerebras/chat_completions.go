// Package cerebras provides the Chat Completions API for the Cerebras Cloud SDK.
//
// The Chat Completions API allows you to generate conversational responses
// using Cerebras Cloud's LLM models with a structured message format.
//
// Basic usage:
//
//	client := cerebras.NewClient(cerebras.WithAPIKey("sk-..."))
//
//	response, err := client.ChatCompletions.Create(ctx,
//		cerebras.Llama31_8b,
//		[]cerebras.Message{
//			{Role: cerebras.MessageRoleUser, Content: "Hello!"},
//		},
//		cerebras.WithTemperature(0.7),
//		cerebras.WithMaxTokens(100),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response.Choices[0].Message.Content)
package cerebras

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	cerebraserrors "github.com/ldaidone/cerebras-cloud-sdk-go/internal/errors"
	"github.com/ldaidone/cerebras-cloud-sdk-go/internal/transport"
)

// ChatCompletionsService provides access to the Chat Completions API.
//
// This service allows you to create chat completions using Cerebras Cloud's
// LLM models. It supports all standard parameters including temperature,
// max_tokens, top_p, and more.
type ChatCompletionsService struct {
	client    *Client
	transport *transport.Transport
}

// ChatCompletionOption is a functional option for chat completion requests.
type ChatCompletionOption func(*ChatCompletionRequest)

// WithTemperature sets the sampling temperature for the completion.
// Higher values (e.g., 0.8) make output more random, lower values (e.g., 0.2) make it more focused.
// Recommended range: 0.0 to 1.5.
func WithTemperature(temp float64) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.Temperature = PtrFloat64(temp)
	}
}

// WithMaxTokens sets the maximum number of tokens to generate.
// This is deprecated in favor of WithMaxCompletionTokens.
func WithMaxTokens(tokens int) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.MaxTokens = PtrInt(tokens)
	}
}

// WithMaxCompletionTokens sets an upper bound for tokens generated,
// including visible output tokens and reasoning tokens.
func WithMaxCompletionTokens(tokens int) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.MaxCompletionTokens = PtrInt(tokens)
	}
}

// WithTopP sets the nucleus sampling probability mass.
// Model considers tokens with top_p probability mass (e.g., 0.1 = top 10%).
func WithTopP(p float64) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.TopP = PtrFloat64(p)
	}
}

// WithStop sets stop sequences where the API will stop generating tokens.
// Up to 4 sequences can be specified.
//
// Performance: Uses strings.Builder for efficient string concatenation.
func WithStop(sequences ...string) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		if len(sequences) == 1 {
			req.Stop = PtrString(sequences[0])
		} else if len(sequences) > 1 {
			// Use strings.Builder for efficient concatenation
			var builder strings.Builder
			for i, s := range sequences {
				if i > 0 {
					builder.WriteString(",")
				}
				builder.WriteString(s)
			}
			req.Stop = PtrString(builder.String())
		}
	}
}

// WithN sets the number of chat completion choices to generate per input message.
func WithN(n int) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.N = PtrInt(n)
	}
}

// WithFrequencyPenalty sets the frequency penalty (-2.0 to 2.0).
// Positive values penalize tokens based on their frequency in the text so far.
func WithFrequencyPenalty(penalty float64) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.FrequencyPenalty = PtrFloat64(penalty)
	}
}

// WithPresencePenalty sets the presence penalty (-2.0 to 2.0).
// Positive values penalize tokens based on whether they appear in the text so far.
func WithPresencePenalty(penalty float64) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.PresencePenalty = PtrFloat64(penalty)
	}
}

// WithLogprobs specifies whether to return log probabilities of output tokens.
func WithLogprobs(logprobs bool) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.Logprobs = PtrBool(logprobs)
	}
}

// WithTopLogprobs sets the number of most likely tokens to return (0-20).
// Requires WithLogprobs(true).
func WithTopLogprobs(n int) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.TopLogprobs = PtrInt(n)
	}
}

// WithLogitBias sets a JSON object mapping token IDs to bias values (-100 to 100).
func WithLogitBias(bias map[string]float64) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.LogitBias = bias
	}
}

// WithSeed sets the seed for deterministic sampling (best-effort).
func WithSeed(seed int) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.Seed = PtrInt(seed)
	}
}

// WithUser sets a unique identifier for end-user monitoring.
func WithUser(user string) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.User = PtrString(user)
	}
}

// WithServiceTier sets the service tier selection.
// Options: "auto", "default", "flex", "priority".
func WithServiceTier(tier string) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ServiceTier = PtrString(tier)
	}
}

// WithResponseFormat sets the response format.
// Use ResponseFormat{Type: "json_object"} for JSON output.
func WithResponseFormat(format ResponseFormat) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ResponseFormat = &format
	}
}

// WithTools sets the list of tools the model may call.
func WithTools(tools ...Tool) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.Tools = tools
	}
}

// WithToolChoice sets the tool choice parameter.
// Can be "none", "auto", "required", or a specific function.
func WithToolChoice(choice interface{}) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ToolChoice = choice
	}
}

// WithParallelToolCalls enables or disables parallel tool calling.
func WithParallelToolCalls(enabled bool) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ParallelToolCalls = PtrBool(enabled)
	}
}

// WithReasoningEffort sets the reasoning effort for reasoning models.
// Options: "low", "medium", "high".
func WithReasoningEffort(effort string) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ReasoningEffort = PtrString(effort)
	}
}

// WithReasoningFormat sets how reasoning is returned.
// Options: "none", "parsed", "text_parsed", "raw", "hidden".
func WithReasoningFormat(format string) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ReasoningFormat = PtrString(format)
	}
}

// WithClearThinking controls whether thinking content from previous turns is included.
// Only supported on zai-glm-4.7 model.
func WithClearThinking(enabled bool) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.ClearThinking = PtrBool(enabled)
	}
}

// WithDisableReasoning disables reasoning for reasoning models.
func WithDisableReasoning(disabled bool) ChatCompletionOption {
	return func(req *ChatCompletionRequest) {
		req.DisableReasoning = PtrBool(disabled)
	}
}

// Create generates a chat completion using the specified model and messages.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - model: The model identifier (e.g., Llama31_8b, Llama31_70b)
//   - messages: The list of messages in the conversation
//   - opts: Optional functional options for additional parameters
//
// Returns:
//   - *ChatCompletion: The chat completion response
//   - error: Any error that occurred during the request
//
// Example:
//
//	response, err := client.ChatCompletions.Create(ctx,
//		cerebras.Llama31_8b,
//		[]cerebras.Message{
//			{Role: cerebras.MessageRoleSystem, Content: "You are a helpful assistant."},
//			{Role: cerebras.MessageRoleUser, Content: "Hello!"},
//		},
//		cerebras.WithTemperature(0.7),
//		cerebras.WithMaxTokens(100),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response.Choices[0].Message.Content)
func (s *ChatCompletionsService) Create(
	ctx context.Context,
	model string,
	messages []Message,
	opts ...ChatCompletionOption,
) (*ChatCompletion, error) {
	// Validate required parameters
	if model == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "model parameter is required"))
	}
	if messages == nil || len(messages) == 0 {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "messages parameter is required and must not be empty"))
	}

	// Build request from parameters and options
	req := &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Make the HTTP request
	var result ChatCompletion
	err := s.client.Do(ctx, http.MethodPost, "/v1/chat/completions", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateWithRequest creates a chat completion using a pre-built ChatCompletionRequest.
//
// This is useful when you need more control over the request or when
// building requests dynamically.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - request: The pre-built chat completion request
//
// Returns:
//   - *ChatCompletion: The chat completion response
//   - error: Any error that occurred during the request
func (s *ChatCompletionsService) CreateWithRequest(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletion, error) {
	// Validate required parameters
	if request.Model == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "model parameter is required"))
	}
	if request.Messages == nil || len(request.Messages) == 0 {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "messages parameter is required and must not be empty"))
	}

	// Make the HTTP request
	var result ChatCompletion
	err := s.client.Do(ctx, http.MethodPost, "/v1/chat/completions", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateStream creates a streaming chat completion.
//
// This method returns two channels: one for receiving stream chunks and one for errors.
// The stream continues until the channel is closed or an error is received.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - model: The model identifier (e.g., Llama31_8b, Llama31_70b)
//   - messages: The list of messages in the conversation
//   - opts: Optional functional options for additional parameters
//
// Returns:
//   - <-chan StreamResponse: Channel receiving stream chunks
//   - <-chan error: Channel receiving errors
//   - error: Error if stream setup fails
//
// Example:
//
//	stream, errChan, err := client.ChatCompletions.CreateStream(ctx,
//		cerebras.Llama31_8b,
//		[]cerebras.Message{
//			{Role: cerebras.MessageRoleUser, Content: "Write a poem"},
//		},
//		cerebras.WithTemperature(0.7),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for {
//		select {
//		case chunk, ok := <-stream:
//			if !ok {
//				return // Stream closed
//			}
//			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
//				fmt.Print(*chunk.Choices[0].Delta.Content)
//			}
//		case err := <-errChan:
//			log.Printf("Stream error: %v", err)
//			return
//		}
//	}
func (s *ChatCompletionsService) CreateStream(
	ctx context.Context,
	model string,
	messages []Message,
	opts ...ChatCompletionOption,
) (<-chan StreamResponse, <-chan error, error) {
	// Build request from parameters and options
	req := &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Set stream=true
	req.Stream = PtrBool(true)

	// Validate required parameters
	if req.Model == "" {
		return nil, nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "model parameter is required"))
	}
	if req.Messages == nil || len(req.Messages) == 0 {
		return nil, nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "messages parameter is required and must not be empty"))
	}

	// Marshal request body
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to marshal request body: %w", err))
	}

	// Create HTTP request
	reqURL := s.transport.GetBaseURL() + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to create request: %w", err))
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	// Execute request with streaming-optimized client
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to execute request: %w", err))
	}

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, cerebraserrors.MapStatusCode(resp.StatusCode, string(body))
	}

	// Create channels
	respChan := make(chan StreamResponse, 10)
	errChan := make(chan error, 1)

	// Start reading stream in goroutine
	go readSSEStream(ctx, resp.Body, respChan, errChan)

	return respChan, errChan, nil
}

// CreateStreamWithRequest creates a streaming chat completion with a pre-built request.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - request: The pre-built chat completion request (Stream will be set to true)
//
// Returns:
//   - <-chan StreamResponse: Channel receiving stream chunks
//   - <-chan error: Channel receiving errors
//   - error: Error if stream setup fails
func (s *ChatCompletionsService) CreateStreamWithRequest(
	ctx context.Context,
	request *ChatCompletionRequest,
) (<-chan StreamResponse, <-chan error, error) {
	// Validate required parameters
	if request.Model == "" {
		return nil, nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "model parameter is required"))
	}
	if request.Messages == nil || len(request.Messages) == 0 {
		return nil, nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "messages parameter is required and must not be empty"))
	}

	// Set stream=true
	request.Stream = PtrBool(true)

	// Marshal request body
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to marshal request body: %w", err))
	}

	// Create HTTP request
	reqURL := s.transport.GetBaseURL() + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to create request: %w", err))
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	// Execute request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to execute request: %w", err))
	}

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, cerebraserrors.MapStatusCode(resp.StatusCode, string(body))
	}

	// Create channels
	respChan := make(chan StreamResponse, 10)
	errChan := make(chan error, 1)

	// Start reading stream in goroutine
	go readSSEStream(ctx, resp.Body, respChan, errChan)

	return respChan, errChan, nil
}
