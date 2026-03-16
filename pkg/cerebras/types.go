// Package cerebras provides types for the Cerebras Cloud SDK.
// This file contains all type definitions for API entities.
package cerebras

// ============================================================================
// Message Types
// ============================================================================

// MessageRole represents the role of a message in a chat conversation.
type MessageRole string

const (
	// MessageRoleSystem represents a system message that sets the behavior of the assistant.
	MessageRoleSystem MessageRole = "system"
	// MessageRoleUser represents a user message in the conversation.
	MessageRoleUser MessageRole = "user"
	// MessageRoleAssistant represents an assistant message in the conversation.
	MessageRoleAssistant MessageRole = "assistant"
)

// Message represents a chat message in a conversation.
type Message struct {
	// Role is the role of the message sender (system, user, or assistant).
	Role MessageRole `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
	// Name is an optional name for the participant (for user messages).
	Name *string `json:"name,omitempty"`
	// ToolCalls is an array of tool calls (for assistant messages).
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	// ToolCallID is the ID of the tool call this message is responding to (for tool responses).
	ToolCallID *string `json:"tool_call_id,omitempty"`
	// Reasoning is the reasoning content for reasoning models.
	Reasoning *string `json:"reasoning,omitempty"`
}

// ============================================================================
// Chat Completion Request Types
// ============================================================================

// ChatCompletionRequest represents a request to the chat completions endpoint.
type ChatCompletionRequest struct {
	// Model is the ID of the model to use (required).
	Model string `json:"model"`
	// Messages is the list of messages comprising the conversation (required).
	Messages []Message `json:"messages"`

	// Optional parameters

	// ClearThinking controls whether thinking content from previous turns is included.
	// Only supported on zai-glm-4.7 model.
	ClearThinking *bool `json:"clear_thinking,omitempty"`
	// DisableReasoning disables reasoning for reasoning models.
	DisableReasoning *bool `json:"disable_reasoning,omitempty"`
	// FrequencyPenalty is a number between -2.0 and 2.0 that penalizes tokens based on their frequency.
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	// LogitBias is a JSON object mapping token IDs to bias values (-100 to 100).
	LogitBias map[string]float64 `json:"logit_bias,omitempty"`
	// Logprobs specifies whether to return log probabilities of output tokens.
	Logprobs *bool `json:"logprobs,omitempty"`
	// MaxCompletionTokens is an upper bound for tokens generated (including reasoning tokens).
	MaxCompletionTokens *int `json:"max_completion_tokens,omitempty"`
	// MaxTokens is the maximum number of tokens to generate (deprecated in favor of MaxCompletionTokens).
	MaxTokens *int `json:"max_tokens,omitempty"`
	// MinCompletionTokens is the minimum number of tokens to generate.
	MinCompletionTokens *int `json:"min_completion_tokens,omitempty"`
	// MinTokens is the minimum number of tokens to generate.
	MinTokens *int `json:"min_tokens,omitempty"`
	// N is the number of chat completion choices to generate per input message.
	N *int `json:"n,omitempty"`
	// ParallelToolCalls enables parallel function calling during tool use.
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`
	// Prediction is the configuration for Predicted Output.
	Prediction *Prediction `json:"prediction,omitempty"`
	// PresencePenalty is a number between -2.0 and 2.0 that penalizes tokens based on whether they appear.
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`
	// ReasoningEffort constrains reasoning effort for reasoning models ("low", "medium", "high").
	ReasoningEffort *string `json:"reasoning_effort,omitempty"`
	// ReasoningFormat determines how reasoning is returned ("none", "parsed", "text_parsed", "raw", "hidden").
	ReasoningFormat *string `json:"reasoning_format,omitempty"`
	// ResponseFormat controls the output format (text, json_object, json_schema).
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	// Seed is for deterministic sampling (best-effort).
	Seed *int `json:"seed,omitempty"`
	// ServiceTier is the service tier selection ("auto", "default", "flex", "priority").
	ServiceTier *string `json:"service_tier,omitempty"`
	// Stop is up to 4 sequences where the API will stop generating further tokens.
	Stop *string `json:"stop,omitempty"`
	// Stream specifies whether to stream the response.
	Stream *bool `json:"stream,omitempty"`
	// StreamOptions contains options for streaming.
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	// Temperature is the sampling temperature between 0 and 1.5.
	Temperature *float64 `json:"temperature,omitempty"`
	// ToolChoice controls which tool the model uses.
	ToolChoice interface{} `json:"tool_choice,omitempty"` // string or ToolChoice object
	// Tools is the list of tools the model may call.
	Tools []Tool `json:"tools,omitempty"`
	// TopLogprobs is the number of most likely tokens to return (0-20).
	TopLogprobs *int `json:"top_logprobs,omitempty"`
	// TopP is the nucleus sampling probability mass.
	TopP *float64 `json:"top_p,omitempty"`
	// User is a unique identifier for end-user monitoring.
	User *string `json:"user,omitempty"`
}

// ============================================================================
// Chat Completion Response Types
// ============================================================================

// ChatCompletion represents a response from the chat completions endpoint.
type ChatCompletion struct {
	// ID is the unique identifier for the chat completion.
	ID string `json:"id"`
	// Object is the type of object (always "chat.completion").
	Object string `json:"object"`
	// Created is the Unix timestamp of creation.
	Created float64 `json:"created"`
	// Model is the model used for the completion.
	Model string `json:"model"`
	// Choices is the list of chat completion choices.
	Choices []Choice `json:"choices"`
	// Usage contains token usage information.
	Usage Usage `json:"usage"`
	// SystemFingerprint is a fingerprint for the system.
	SystemFingerprint string `json:"system_fingerprint"`
	// ServiceTierUsed is the tier used (only when service_tier: auto).
	ServiceTierUsed *string `json:"service_tier_used,omitempty"`
	// TimeInfo contains timing information.
	TimeInfo TimeInfo `json:"time_info"`
}

// Choice represents a chat completion choice.
type Choice struct {
	// FinishReason is the reason the model stopped generating tokens.
	FinishReason FinishReason `json:"finish_reason"`
	// Index is the index of the choice in the list.
	Index int `json:"index"`
	// Message is the assistant message.
	Message Message `json:"message"`
	// Logprobs contains log probability information (if requested).
	Logprobs *Logprobs `json:"logprobs,omitempty"`
}

// FinishReason represents the reason the model stopped generating tokens.
type FinishReason string

const (
	// FinishReasonStop means the model hit a natural stop point or a provided stop sequence.
	FinishReasonStop FinishReason = "stop"
	// FinishReasonLength means the maximum number of tokens specified was reached.
	FinishReasonLength FinishReason = "length"
	// FinishReasonContentFilter means content was omitted due to a flag from content filters.
	FinishReasonContentFilter FinishReason = "content_filter"
	// FinishReasonToolCalls means the model called a tool.
	FinishReasonToolCalls FinishReason = "tool_calls"
)

// Usage represents token usage information.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`
	// PromptTokensDetails contains details about prompt tokens.
	PromptTokensDetails *PromptTokensDetails `json:"prompt_tokens_details,omitempty"`
	// CompletionTokensDetails contains details about completion tokens.
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

// PromptTokensDetails contains details about prompt tokens.
type PromptTokensDetails struct {
	// CachedTokens is the number of tokens served from cache.
	CachedTokens int `json:"cached_tokens"`
}

// CompletionTokensDetails contains details about completion tokens.
type CompletionTokensDetails struct {
	// AcceptedPredictionTokens is the number of prediction tokens that appeared in completion.
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
	// RejectedPredictionTokens is the number of prediction tokens not in completion.
	RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
}

// TimeInfo represents timing information for the completion.
type TimeInfo struct {
	// QueueTime is the time spent in queue (seconds).
	QueueTime float64 `json:"queue_time"`
	// PromptTime is the time processing the prompt (seconds).
	PromptTime float64 `json:"prompt_time"`
	// CompletionTime is the time generating the completion (seconds).
	CompletionTime float64 `json:"completion_time"`
	// TotalTime is the total request time (seconds).
	TotalTime float64 `json:"total_time"`
	// Created is the Unix timestamp of creation.
	Created float64 `json:"created"`
}

// Logprobs represents log probability information for a choice.
type Logprobs struct {
	// Content is the list of token log probabilities.
	Content []TokenLogprob `json:"content"`
}

// TokenLogprob represents log probability for a single token.
type TokenLogprob struct {
	// Token is the token text.
	Token string `json:"token"`
	// Logprob is the log probability of the token.
	Logprob float64 `json:"logprob"`
	// Bytes is the byte representation of the token (if available).
	Bytes []int `json:"bytes,omitempty"`
	// TopLogprobs is the list of most likely tokens and their log probabilities.
	TopLogprobs []TopLogprob `json:"top_logprobs,omitempty"`
}

// TopLogprob represents a token and its log probability in the top candidates.
type TopLogprob struct {
	// Token is the token text.
	Token string `json:"token"`
	// Logprob is the log probability of the token.
	Logprob float64 `json:"logprob"`
	// Bytes is the byte representation of the token (if available).
	Bytes []int `json:"bytes,omitempty"`
}

// ============================================================================
// Model Types
// ============================================================================

// Model represents a model available in the Cerebras Cloud API.
type Model struct {
	// ID is the model identifier.
	ID string `json:"id"`
	// Object is the type of object (always "model").
	Object string `json:"object"`
	// Created is the Unix timestamp of creation.
	Created float64 `json:"created"`
	// OwnedBy is the organization that owns the model.
	OwnedBy string `json:"owned_by"`
}

// ModelList represents a list of models.
type ModelList struct {
	// Object is the type of object (always "list").
	Object string `json:"object"`
	// Data is the list of models.
	Data []Model `json:"data"`
}

// ============================================================================
// Streaming Types
// ============================================================================

// StreamResponse represents a streaming response chunk.
type StreamResponse struct {
	// ID is the unique identifier for the chat completion.
	ID string `json:"id"`
	// Object is the type of object (always "chat.completion.chunk").
	Object string `json:"object"`
	// Created is the Unix timestamp of creation.
	Created float64 `json:"created"`
	// Model is the model used for the completion.
	Model string `json:"model"`
	// Choices is the list of streaming choices.
	Choices []StreamChoice `json:"choices"`
	// Usage contains token usage (only in final chunk if include_usage is set).
	Usage *Usage `json:"usage,omitempty"`
	// TimeInfo contains timing information.
	TimeInfo *TimeInfo `json:"time_info,omitempty"`
}

// StreamChoice represents a streaming choice with delta content.
type StreamChoice struct {
	// Delta is the incremental content update.
	Delta Delta `json:"delta"`
	// FinishReason is the reason the model stopped (only in final chunk).
	FinishReason *FinishReason `json:"finish_reason,omitempty"`
	// Index is the index of the choice.
	Index int `json:"index"`
	// Logprobs contains log probability information (if requested).
	Logprobs *Logprobs `json:"logprobs,omitempty"`
}

// Delta represents incremental content update in streaming responses.
type Delta struct {
	// Role is the role of the message (only in first chunk).
	Role *MessageRole `json:"role,omitempty"`
	// Content is the incremental content update.
	Content *string `json:"content,omitempty"`
	// Reasoning is the reasoning content update (for reasoning models).
	Reasoning *string `json:"reasoning,omitempty"`
	// ToolCalls is the list of tool calls (if applicable).
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// StreamOptions contains options for streaming responses.
type StreamOptions struct {
	// IncludeUsage specifies whether to include usage in the final chunk.
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

// ============================================================================
// Tool Calling Types
// ============================================================================

// Tool represents a tool (function) that the model may call.
type Tool struct {
	// Type is the type of tool (always "function" for now).
	Type string `json:"type"`
	// Function is the function definition.
	Function Function `json:"function"`
}

// Function represents a function definition for tool calling.
type Function struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is a description of what the function does.
	Description string `json:"description"`
	// Parameters is the JSON Schema for the function parameters.
	Parameters interface{} `json:"parameters"` // JSON Schema object
	// Strict specifies whether to enforce strict JSON Schema validation.
	Strict *bool `json:"strict,omitempty"`
}

// ToolChoice controls which tool the model uses.
// Can be a string ("none", "auto", "required") or a ToolChoiceObject.
type ToolChoice struct {
	// Type is the type of tool choice ("function").
	Type string `json:"type"`
	// Function is the function to call.
	Function ToolChoiceFunction `json:"function"`
}

// ToolChoiceFunction represents a function reference for tool choice.
type ToolChoiceFunction struct {
	// Name is the name of the function to call.
	Name string `json:"name"`
}

// ToolCall represents a tool call made by the model.
type ToolCall struct {
	// ID is the unique identifier for the tool call.
	ID string `json:"id"`
	// Type is the type of tool call ("function").
	Type string `json:"type"`
	// Function is the function that was called.
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents a function call in a tool call.
type ToolCallFunction struct {
	// Name is the name of the function called.
	Name string `json:"name"`
	// Arguments is the JSON-encoded arguments string.
	Arguments string `json:"arguments"`
}

// ============================================================================
// Response Format Types
// ============================================================================

// ResponseFormat controls the output format of the completion.
type ResponseFormat struct {
	// Type is the format type ("text", "json_object", "json_schema").
	Type string `json:"type"`
	// JSONSchema is the schema definition (only for json_schema type).
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

// JSONSchema represents a JSON Schema for structured output.
type JSONSchema struct {
	// Name is the name of the schema.
	Name string `json:"name"`
	// Description is a description of the schema.
	Description string `json:"description"`
	// Schema is the JSON Schema definition.
	Schema interface{} `json:"schema"` // JSON Schema object
	// Strict specifies whether to enforce strict JSON Schema validation.
	Strict *bool `json:"strict,omitempty"`
}

// ============================================================================
// Additional Types
// ============================================================================

// Prediction represents configuration for Predicted Output.
type Prediction struct {
	// Type is the prediction type.
	Type string `json:"type"`
	// Content is the predicted content.
	Content []PredictionContent `json:"content"`
}

// PredictionContent represents predicted content.
type PredictionContent struct {
	// Type is the content type ("text").
	Type string `json:"type"`
	// Text is the predicted text.
	Text string `json:"text"`
}

// ============================================================================
// Model Constants
// ============================================================================

// Model identifiers for use with ChatCompletionRequest.Model
const (
	// Llama31_8b is the Llama 3.1 8B model from Meta.
	Llama31_8b = "llama3.1-8b"
	// Llama31_70b is the Llama 3.1 70B model from Meta.
	Llama31_70b = "llama3.1-70b"
	// GptOss120b is the GPT-oss 120B model.
	GptOss120b = "gpt-oss-120b"
	// Qwen3_235b is the Qwen 3 235B A22B model.
	Qwen3_235b = "qwen-3-235b-a22b"
	// ZaiGlm47 is the ZAI GLM 4.7 model.
	ZaiGlm47 = "zai-glm-4.7"
)

// ============================================================================
// Helper Functions
// ============================================================================

// PtrString returns a pointer to a string value.
// Useful for setting optional string parameters in requests.
func PtrString(s string) *string {
	return &s
}

// PtrInt returns a pointer to an int value.
// Useful for setting optional int parameters in requests.
func PtrInt(i int) *int {
	return &i
}

// PtrFloat64 returns a pointer to a float64 value.
// Useful for setting optional float parameters in requests.
func PtrFloat64(f float64) *float64 {
	return &f
}

// PtrBool returns a pointer to a bool value.
// Useful for setting optional bool parameters in requests.
func PtrBool(b bool) *bool {
	return &b
}
