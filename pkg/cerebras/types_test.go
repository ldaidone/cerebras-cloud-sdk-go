package cerebras

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageRole(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected string
	}{
		{"system role", MessageRoleSystem, "system"},
		{"user role", MessageRoleUser, "user"},
		{"assistant role", MessageRoleAssistant, "assistant"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("MessageRole = %s, want %s", tt.role, tt.expected)
			}
		})
	}
}

func TestMessage_JSON(t *testing.T) {
	name := "test-user"
	msg := Message{
		Role:    MessageRoleUser,
		Content: "Hello, world!",
		Name:    &name,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal Message: %v", err)
	}

	expected := `{"role":"user","content":"Hello, world!","name":"test-user"}`
	if string(data) != expected {
		t.Errorf("Marshaled JSON = %s, want %s", string(data), expected)
	}

	var unmarshaled Message
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Message: %v", err)
	}

	if unmarshaled.Role != msg.Role {
		t.Errorf("Role = %s, want %s", unmarshaled.Role, msg.Role)
	}
	if unmarshaled.Content != msg.Content {
		t.Errorf("Content = %s, want %s", unmarshaled.Content, msg.Content)
	}
}

func TestChatCompletionRequest_JSON(t *testing.T) {
	temp := 0.7
	maxTokens := 100
	req := ChatCompletionRequest{
		Model: Llama31_8b,
		Messages: []Message{
			{Role: MessageRoleUser, Content: "Hello!"},
		},
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatCompletionRequest: %v", err)
	}

	// Verify snake_case field names
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if _, ok := jsonMap["max_tokens"]; !ok {
		t.Error("Expected 'max_tokens' field (snake_case)")
	}
	if _, ok := jsonMap["temperature"]; !ok {
		t.Error("Expected 'temperature' field (snake_case)")
	}
}

func TestChatCompletion_JSON(t *testing.T) {
	completion := ChatCompletion{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   Llama31_8b,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    MessageRoleAssistant,
					Content: "Hello! How can I help?",
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

	data, err := json.Marshal(completion)
	if err != nil {
		t.Fatalf("Failed to marshal ChatCompletion: %v", err)
	}

	var unmarshaled ChatCompletion
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ChatCompletion: %v", err)
	}

	if unmarshaled.ID != completion.ID {
		t.Errorf("ID = %s, want %s", unmarshaled.ID, completion.ID)
	}
	if len(unmarshaled.Choices) != 1 {
		t.Errorf("Choices length = %d, want 1", len(unmarshaled.Choices))
	}
}

func TestFinishReason(t *testing.T) {
	tests := []struct {
		name     string
		reason   FinishReason
		expected string
	}{
		{"stop", FinishReasonStop, "stop"},
		{"length", FinishReasonLength, "length"},
		{"content_filter", FinishReasonContentFilter, "content_filter"},
		{"tool_calls", FinishReasonToolCalls, "tool_calls"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.reason) != tt.expected {
				t.Errorf("FinishReason = %s, want %s", tt.reason, tt.expected)
			}
		})
	}
}

func TestUsage_JSON(t *testing.T) {
	usage := Usage{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		PromptTokensDetails: &PromptTokensDetails{
			CachedTokens: 80,
		},
		CompletionTokensDetails: &CompletionTokensDetails{
			AcceptedPredictionTokens: 10,
			RejectedPredictionTokens: 5,
		},
	}

	data, err := json.Marshal(usage)
	if err != nil {
		t.Fatalf("Failed to marshal Usage: %v", err)
	}

	var unmarshaled Usage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Usage: %v", err)
	}

	if unmarshaled.TotalTokens != usage.TotalTokens {
		t.Errorf("TotalTokens = %d, want %d", unmarshaled.TotalTokens, usage.TotalTokens)
	}
}

func TestModel_JSON(t *testing.T) {
	model := Model{
		ID:      Llama31_8b,
		Object:  "model",
		Created: time.Now().Unix(),
		OwnedBy: "meta",
	}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("Failed to marshal Model: %v", err)
	}

	var unmarshaled Model
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Model: %v", err)
	}

	if unmarshaled.ID != model.ID {
		t.Errorf("ID = %s, want %s", unmarshaled.ID, model.ID)
	}
}

func TestStreamResponse_JSON(t *testing.T) {
	content := "Hello"
	role := MessageRoleAssistant
	streamResp := StreamResponse{
		ID:      "test-id",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   Llama31_8b,
		Choices: []StreamChoice{
			{
				Index: 0,
				Delta: Delta{
					Role:    &role,
					Content: &content,
				},
			},
		},
	}

	data, err := json.Marshal(streamResp)
	if err != nil {
		t.Fatalf("Failed to marshal StreamResponse: %v", err)
	}

	var unmarshaled StreamResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal StreamResponse: %v", err)
	}

	if len(unmarshaled.Choices) != 1 {
		t.Errorf("Choices length = %d, want 1", len(unmarshaled.Choices))
	}
}

func TestTool_JSON(t *testing.T) {
	params := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city and state",
			},
		},
	}

	tool := Tool{
		Type: "function",
		Function: Function{
			Name:        "get_weather",
			Description: "Get the weather in a location",
			Parameters:  params,
		},
	}

	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("Failed to marshal Tool: %v", err)
	}

	var unmarshaled Tool
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Tool: %v", err)
	}

	if unmarshaled.Function.Name != "get_weather" {
		t.Errorf("Function.Name = %s, want get_weather", unmarshaled.Function.Name)
	}
}

func TestResponseFormat_JSON(t *testing.T) {
	rf := ResponseFormat{
		Type: "json_object",
	}

	data, err := json.Marshal(rf)
	if err != nil {
		t.Fatalf("Failed to marshal ResponseFormat: %v", err)
	}

	var unmarshaled ResponseFormat
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ResponseFormat: %v", err)
	}

	if unmarshaled.Type != "json_object" {
		t.Errorf("Type = %s, want json_object", unmarshaled.Type)
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("PtrString", func(t *testing.T) {
		s := "test"
		ptr := PtrString(s)
		if *ptr != s {
			t.Errorf("*ptr = %s, want %s", *ptr, s)
		}
	})

	t.Run("PtrInt", func(t *testing.T) {
		i := 42
		ptr := PtrInt(i)
		if *ptr != i {
			t.Errorf("*ptr = %d, want %d", *ptr, i)
		}
	})

	t.Run("PtrFloat64", func(t *testing.T) {
		f := 3.14
		ptr := PtrFloat64(f)
		if *ptr != f {
			t.Errorf("*ptr = %f, want %f", *ptr, f)
		}
	})

	t.Run("PtrBool", func(t *testing.T) {
		b := true
		ptr := PtrBool(b)
		if *ptr != b {
			t.Errorf("*ptr = %v, want %v", *ptr, b)
		}
	})
}

func TestModelConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Llama31_8b", Llama31_8b, "llama3.1-8b"},
		{"Llama31_70b", Llama31_70b, "llama3.1-70b"},
		{"GptOss120b", GptOss120b, "gpt-oss-120b"},
		{"Qwen3_235b", Qwen3_235b, "qwen-3-235b-a22b"},
		{"ZaiGlm47", ZaiGlm47, "zai-glm-4.7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant = %s, want %s", tt.constant, tt.expected)
			}
		})
	}
}
