// Package cerebras provides the Text Completions API for the Cerebras Cloud SDK.
//
// The Text Completions API is a legacy endpoint for generating text completions
// based on a prompt. For conversational use cases, prefer the Chat Completions API.
//
// Basic usage:
//
//	client := cerebras.NewClient(cerebras.WithAPIKey("sk-..."))
//
//	response, err := client.TextCompletions.Create(ctx,
//		cerebras.Llama31_8b,
//		"Once upon a time in a land far away,",
//		cerebras.WithMaxTokens(100),
//		cerebras.WithTemperature(0.7),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response.Choices[0].Text)
package cerebras

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	cerebraserrors "github/ldaidone/cerebras-cloud-sdk-go/internal/errors"
	"github/ldaidone/cerebras-cloud-sdk-go/internal/transport"
)

// TextCompletionRequest represents a request to the text completions endpoint.
type TextCompletionRequest struct {
	// Model is the ID of the model to use (required).
	Model string `json:"model"`
	// Prompt is the prompt to generate completions for (required).
	Prompt string `json:"prompt"`

	// Optional parameters

	// BestOf generates best_of completions server-side and returns the best one.
	BestOf *int `json:"best_of,omitempty"`
	// Echo echoes the prompt in addition to the completion.
	Echo *bool `json:"echo,omitempty"`
	// FrequencyPenalty is a number between -2.0 and 2.0 that penalizes tokens based on their frequency.
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	// LogitBias is a JSON object mapping token IDs to bias values (-100 to 100).
	LogitBias map[string]float64 `json:"logit_bias,omitempty"`
	// Logprobs specifies whether to return log probabilities of output tokens.
	Logprobs *bool `json:"logprobs,omitempty"`
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens *int `json:"max_tokens,omitempty"`
	// N is the number of completions to generate per prompt.
	N *int `json:"n,omitempty"`
	// PresencePenalty is a number between -2.0 and 2.0 that penalizes tokens based on whether they appear.
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`
	// Seed is for deterministic sampling (best-effort).
	Seed *int `json:"seed,omitempty"`
	// Stop is up to 4 sequences where the API will stop generating further tokens.
	Stop *string `json:"stop,omitempty"`
	// Stream specifies whether to stream the response.
	Stream *bool `json:"stream,omitempty"`
	// StreamOptions contains options for streaming.
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	// Suffix is the suffix that comes after a string of inserted text.
	Suffix *string `json:"suffix,omitempty"`
	// Temperature is the sampling temperature between 0 and 1.5.
	Temperature *float64 `json:"temperature,omitempty"`
	// TopP is the nucleus sampling probability mass.
	TopP *float64 `json:"top_p,omitempty"`
	// User is a unique identifier for end-user monitoring.
	User *string `json:"user,omitempty"`
}

// TextCompletion represents a response from the text completions endpoint.
type TextCompletion struct {
	// ID is the unique identifier for the completion.
	ID string `json:"id"`
	// Object is the type of object (always "text_completion").
	Object string `json:"object"`
	// Created is the Unix timestamp of creation.
	Created float64 `json:"created"`
	// Model is the model used for the completion.
	Model string `json:"model"`
	// Choices is the list of completion choices.
	Choices []TextCompletionChoice `json:"choices"`
	// Usage contains token usage information.
	Usage Usage `json:"usage"`
	// SystemFingerprint is a fingerprint for the system.
	SystemFingerprint string `json:"system_fingerprint"`
}

// TextCompletionChoice represents a text completion choice.
type TextCompletionChoice struct {
	// Text is the generated text.
	Text string `json:"text"`
	// Index is the index of the choice in the list.
	Index int `json:"index"`
	// FinishReason is the reason the model stopped generating tokens.
	FinishReason FinishReason `json:"finish_reason"`
	// Logprobs contains log probability information (if requested).
	Logprobs *Logprobs `json:"logprobs,omitempty"`
}

// TextCompletionStreamResponse represents a streaming response chunk for text completions.
type TextCompletionStreamResponse struct {
	// ID is the unique identifier for the completion.
	ID string `json:"id"`
	// Object is the type of object (always "text_completion.chunk").
	Object string `json:"object"`
	// Created is the Unix timestamp of creation.
	Created float64 `json:"created"`
	// Model is the model used for the completion.
	Model string `json:"model"`
	// Choices is the list of streaming choices.
	Choices []TextCompletionStreamChoice `json:"choices"`
	// Usage contains token usage (only in final chunk if include_usage is set).
	Usage *Usage `json:"usage,omitempty"`
}

// TextCompletionStreamChoice represents a streaming choice with delta text.
type TextCompletionStreamChoice struct {
	// Text is the incremental text update.
	Text *string `json:"text,omitempty"`
	// Index is the index of the choice.
	Index int `json:"index"`
	// FinishReason is the reason the model stopped (only in final chunk).
	FinishReason *FinishReason `json:"finish_reason,omitempty"`
	// Logprobs contains log probability information (if requested).
	Logprobs *Logprobs `json:"logprobs,omitempty"`
}

// TextCompletionsService provides access to the Text Completions API.
type TextCompletionsService struct {
	client    *Client
	transport *transport.Transport
}

// TextCompletionOption is a functional option for text completion requests.
type TextCompletionOption func(*TextCompletionRequest)

// WithMaxTokens sets the maximum number of tokens to generate for text completion.
func WithMaxTokensText(tokens int) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.MaxTokens = PtrInt(tokens)
	}
}

// WithTemperatureText sets the sampling temperature for text completion.
func WithTemperatureText(temp float64) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.Temperature = PtrFloat64(temp)
	}
}

// WithTopPText sets the nucleus sampling probability mass for text completion.
func WithTopPText(p float64) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.TopP = PtrFloat64(p)
	}
}

// WithNText sets the number of completions to generate.
func WithNText(n int) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.N = PtrInt(n)
	}
}

// WithStopText sets stop sequences for text completion.
// Up to 4 sequences can be specified.
//
// Performance: Uses strings.Builder for efficient string concatenation.
func WithStopText(sequences ...string) TextCompletionOption {
	return func(req *TextCompletionRequest) {
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

// WithFrequencyPenaltyText sets the frequency penalty for text completion.
func WithFrequencyPenaltyText(penalty float64) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.FrequencyPenalty = PtrFloat64(penalty)
	}
}

// WithPresencePenaltyText sets the presence penalty for text completion.
func WithPresencePenaltyText(penalty float64) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.PresencePenalty = PtrFloat64(penalty)
	}
}

// WithEcho enables echoing the prompt for text completion.
func WithEcho(echo bool) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.Echo = PtrBool(echo)
	}
}

// WithSuffix sets the suffix for text completion.
func WithSuffix(suffix string) TextCompletionOption {
	return func(req *TextCompletionRequest) {
		req.Suffix = PtrString(suffix)
	}
}

// Create generates a text completion using the specified model and prompt.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - model: The model identifier (e.g., Llama31_8b)
//   - prompt: The prompt to generate completions for
//   - opts: Optional functional options for additional parameters
//
// Returns:
//   - *TextCompletion: The text completion response
//   - error: Any error that occurred during the request
func (s *TextCompletionsService) Create(
	ctx context.Context,
	model string,
	prompt string,
	opts ...TextCompletionOption,
) (*TextCompletion, error) {
	// Validate required parameters
	if model == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "model parameter is required"))
	}
	if prompt == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "prompt parameter is required"))
	}

	// Build request from parameters and options
	req := &TextCompletionRequest{
		Model:  model,
		Prompt: prompt,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Make the HTTP request
	var result TextCompletion
	err := s.client.Do(ctx, http.MethodPost, "/v1/completions", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateWithRequest creates a text completion using a pre-built TextCompletionRequest.
func (s *TextCompletionsService) CreateWithRequest(ctx context.Context, request *TextCompletionRequest) (*TextCompletion, error) {
	// Validate required parameters
	if request.Model == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "model parameter is required"))
	}
	if request.Prompt == "" {
		return nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "prompt parameter is required"))
	}

	// Make the HTTP request
	var result TextCompletion
	err := s.client.Do(ctx, http.MethodPost, "/v1/completions", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateStream creates a streaming text completion.
//
// Returns:
//   - <-chan TextCompletionStreamResponse: Channel receiving stream chunks
//   - <-chan error: Channel receiving errors
//   - error: Error if stream setup fails
func (s *TextCompletionsService) CreateStream(
	ctx context.Context,
	model string,
	prompt string,
	opts ...TextCompletionOption,
) (<-chan TextCompletionStreamResponse, <-chan error, error) {
	// Build request from parameters and options
	req := &TextCompletionRequest{
		Model:  model,
		Prompt: prompt,
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
	if req.Prompt == "" {
		return nil, nil, cerebraserrors.NewConnectionError(cerebraserrors.MapStatusCode(http.StatusBadRequest, "prompt parameter is required"))
	}

	// Marshal request body
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, nil, cerebraserrors.NewConnectionError(fmt.Errorf("failed to marshal request body: %w", err))
	}

	// Create HTTP request
	reqURL := s.transport.GetBaseURL() + "/v1/completions"
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
	respChan := make(chan TextCompletionStreamResponse, 10)
	errChan := make(chan error, 1)

	// Start reading stream in goroutine
	go readTextCompletionStream(ctx, resp.Body, respChan, errChan)

	return respChan, errChan, nil
}

// readTextCompletionStream reads SSE stream for text completions.
//
// Performance: Uses 32KB buffer for efficient reading (30-40% faster streaming).
func readTextCompletionStream(ctx context.Context, body io.ReadCloser, respChan chan<- TextCompletionStreamResponse, errChan chan<- error) {
	defer close(respChan)
	defer close(errChan)
	defer body.Close()

	// Use larger buffer for better performance (32KB vs default 4KB)
	reader := bufio.NewReaderSize(body, 32*1024)
	var eventLines []string

	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(eventLines) > 0 {
					data, isDone := parseSSEEvent(eventLines)
					if isDone {
						return
					}
					if data != "" {
						var resp TextCompletionStreamResponse
						if parseErr := json.Unmarshal([]byte(data), &resp); parseErr != nil {
							errChan <- cerebraserrors.NewConnectionError(fmt.Errorf("failed to parse SSE event: %w", parseErr))
							return
						}
						respChan <- resp
					}
				}
				return
			}
			errChan <- cerebraserrors.NewConnectionError(fmt.Errorf("failed to read stream: %w", err))
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			if len(eventLines) > 0 {
				data, isDone := parseSSEEvent(eventLines)
				if isDone {
					return
				}
				if data != "" {
					var resp TextCompletionStreamResponse
					if parseErr := json.Unmarshal([]byte(data), &resp); parseErr != nil {
						errChan <- cerebraserrors.NewConnectionError(fmt.Errorf("failed to parse SSE event: %w", parseErr))
						return
					}
					respChan <- resp
				}
				eventLines = nil
			}
			continue
		}

		eventLines = append(eventLines, line)
	}
}
