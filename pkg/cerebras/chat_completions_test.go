package cerebras

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChatCompletionsService_Create(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("Path = %s, want /v1/chat/completions", r.URL.Path)
		}

		// Verify headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("Authorization = %s, want Bearer test-key", auth)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", ct)
		}

		// Read request body
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify request parameters
		if req.Model != Llama31_8b {
			t.Errorf("Model = %s, want %s", req.Model, Llama31_8b)
		}
		if len(req.Messages) != 1 {
			t.Errorf("Messages length = %d, want 1", len(req.Messages))
		}
		if req.Temperature == nil || *req.Temperature != 0.7 {
			t.Errorf("Temperature = %v, want 0.7", req.Temperature)
		}

		// Send response
		response := ChatCompletion{
			ID:      "test-completion-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   Llama31_8b,
			Choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    MessageRoleAssistant,
						Content: "Hello! How can I help you today?",
					},
					FinishReason: FinishReasonStop,
				},
			},
			Usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
			TimeInfo: TimeInfo{
				QueueTime:      0.001,
				PromptTime:     0.002,
				CompletionTime: 0.005,
				TotalTime:      0.008,
				Created:        time.Now().Unix(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithTCPWarming(false),
	)

	// Test Create method
	ctx := context.Background()
	response, err := client.ChatCompletions.Create(
		ctx,
		Llama31_8b,
		[]Message{
			{Role: MessageRoleUser, Content: "Hello!"},
		},
		WithTemperature(0.7),
		WithMaxTokens(100),
	)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if response.ID != "test-completion-id" {
		t.Errorf("ID = %s, want test-completion-id", response.ID)
	}
	if len(response.Choices) != 1 {
		t.Errorf("Choices length = %d, want 1", len(response.Choices))
	}
	if response.Choices[0].Message.Content != "Hello! How can I help you today?" {
		t.Errorf("Content = %s, want expected response", response.Choices[0].Message.Content)
	}
	if response.Usage.TotalTokens != 30 {
		t.Errorf("TotalTokens = %d, want 30", response.Usage.TotalTokens)
	}
}

func TestChatCompletionsService_Create_Validation(t *testing.T) {
	client := NewClient(
		WithAPIKey("test-key"),
		WithTCPWarming(false),
	)

	ctx := context.Background()

	t.Run("missing model", func(t *testing.T) {
		_, err := client.ChatCompletions.Create(ctx, "", []Message{{Role: MessageRoleUser, Content: "test"}})
		if err == nil {
			t.Error("Create() should return error for missing model")
		}
	})

	t.Run("missing messages", func(t *testing.T) {
		_, err := client.ChatCompletions.Create(ctx, Llama31_8b, nil)
		if err == nil {
			t.Error("Create() should return error for missing messages")
		}
	})

	t.Run("empty messages array", func(t *testing.T) {
		_, err := client.ChatCompletions.Create(ctx, Llama31_8b, []Message{})
		if err == nil {
			t.Error("Create() should return error for empty messages array")
		}
	})
}

func TestChatCompletionsService_CreateWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify optional parameters
		if req.MaxTokens == nil || *req.MaxTokens != 150 {
			t.Errorf("MaxTokens = %v, want 150", req.MaxTokens)
		}
		if req.TopP == nil || *req.TopP != 0.9 {
			t.Errorf("TopP = %v, want 0.9", req.TopP)
		}
		if req.Stop == nil || *req.Stop != "STOP" {
			t.Errorf("Stop = %v, want STOP", req.Stop)
		}
		if req.N == nil || *req.N != 2 {
			t.Errorf("N = %v, want 2", req.N)
		}
		if req.FrequencyPenalty == nil || *req.FrequencyPenalty != 0.5 {
			t.Errorf("FrequencyPenalty = %v, want 0.5", req.FrequencyPenalty)
		}
		if req.PresencePenalty == nil || *req.PresencePenalty != 0.5 {
			t.Errorf("PresencePenalty = %v, want 0.5", req.PresencePenalty)
		}
		if req.Seed == nil || *req.Seed != 42 {
			t.Errorf("Seed = %v, want 42", req.Seed)
		}

		response := ChatCompletion{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   Llama31_8b,
			Choices: []Choice{{Index: 0, Message: Message{Role: MessageRoleAssistant, Content: "test"}, FinishReason: FinishReasonStop}},
			Usage:   Usage{PromptTokens: 10, CompletionTokens: 10, TotalTokens: 20},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithTCPWarming(false),
	)

	ctx := context.Background()
	_, err := client.ChatCompletions.Create(
		ctx,
		Llama31_8b,
		[]Message{{Role: MessageRoleUser, Content: "test"}},
		WithMaxTokens(150),
		WithTopP(0.9),
		WithStop("STOP"),
		WithN(2),
		WithFrequencyPenalty(0.5),
		WithPresencePenalty(0.5),
		WithSeed(42),
	)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestChatCompletionsService_Create_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithTCPWarming(false),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.ChatCompletions.Create(ctx, Llama31_8b, []Message{{Role: MessageRoleUser, Content: "test"}})
	if err == nil {
		t.Fatal("Create() should return error for cancelled context")
	}
}

func TestChatCompletionsService_CreateWithRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatCompletion{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   Llama31_8b,
			Choices: []Choice{{Index: 0, Message: Message{Role: MessageRoleAssistant, Content: "test"}, FinishReason: FinishReasonStop}},
			Usage:   Usage{PromptTokens: 10, CompletionTokens: 10, TotalTokens: 20},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithTCPWarming(false),
	)

	ctx := context.Background()
	req := &ChatCompletionRequest{
		Model: Llama31_8b,
		Messages: []Message{
			{Role: MessageRoleUser, Content: "test"},
		},
		Temperature: PtrFloat64(0.7),
	}

	response, err := client.ChatCompletions.CreateWithRequest(ctx, req)
	if err != nil {
		t.Fatalf("CreateWithRequest() error = %v", err)
	}

	if response.ID != "test-id" {
		t.Errorf("ID = %s, want test-id", response.ID)
	}
}

func TestChatCompletionsService_ErrorHandling(t *testing.T) {
	t.Run("401 Unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"message": "Invalid API key"}}`))
		}))
		defer server.Close()

		client := NewClient(
			WithAPIKey("invalid-key"),
			WithBaseURL(server.URL),
			WithTCPWarming(false),
		)

		_, err := client.ChatCompletions.Create(context.Background(), Llama31_8b, []Message{{Role: MessageRoleUser, Content: "test"}})
		if err == nil {
			t.Fatal("Create() should return error for 401")
		}
	})

	t.Run("429 Rate Limit", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": {"message": "Rate limit exceeded"}}`))
		}))
		defer server.Close()

		client := NewClient(
			WithAPIKey("test-key"),
			WithBaseURL(server.URL),
			WithTCPWarming(false),
			WithMaxRetries(0), // Disable retries for this test
		)

		_, err := client.ChatCompletions.Create(context.Background(), Llama31_8b, []Message{{Role: MessageRoleUser, Content: "test"}})
		if err == nil {
			t.Fatal("Create() should return error for 429")
		}
	})
}
