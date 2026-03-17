// Package transport benchmarks for performance regression testing.
//
// Run benchmarks with:
//   go test -bench=. -benchmem ./internal/transport/
//
// Compare results to track performance improvements:
//   - Allocation reduction target: 40-50%
//   - JSON processing target: 20-30% faster
//   - Streaming target: 30-40% faster
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a chat completion request.
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
}

// ChatCompletion represents a chat completion response.
type ChatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created float64 `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// BenchmarkMarshal benchmarks JSON marshaling with pooled buffers.
// Expected: 20-30% faster than standard json.Marshal
func BenchmarkMarshal(b *testing.B) {
	req := ChatCompletionRequest{
		Model: "llama3.1-8b",
		Messages: []Message{
			{Role: "user", Content: "Hello, world!"},
		},
		Temperature: float64Ptr(0.7),
		MaxTokens:   intPtr(100),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := getBuffer()
		encoder := json.NewEncoder(buf)
		_ = encoder.Encode(req)
		putBuffer(buf)
	}
}

// BenchmarkUnmarshal benchmarks JSON unmarshaling with Decoder.
// Expected: 20-30% faster than standard json.Unmarshal
func BenchmarkUnmarshal(b *testing.B) {
	jsonData := []byte(`{
		"id": "test-123",
		"object": "chat.completion",
		"created": 1234567890.123,
		"model": "llama3.1-8b",
		"choices": [{
			"index": 0,
			"message": {"role": "assistant", "content": "Hello!"},
			"finish_reason": "stop"
		}],
		"usage": {"prompt_tokens": 10, "completion_tokens": 20, "total_tokens": 30},
		"time_info": {
			"queue_time": 0.001,
			"prompt_time": 0.002,
			"completion_time": 0.005,
			"total_time": 0.008,
			"created": 1234567890.123
		}
	}`)

	var result ChatCompletion

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := json.NewDecoder(bytes.NewReader(jsonData))
		_ = decoder.Decode(&result)
	}
}

// BenchmarkDoRequest benchmarks end-to-end request with pooling.
// Expected: 40-50% allocation reduction
func BenchmarkDoRequest(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      "test-123",
			"object":  "chat.completion",
			"created": 1234567890.123,
			"model":   "llama3.1-8b",
			"choices": []interface{}{
				map[string]interface{}{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Hello!",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 20,
				"total_tokens":      30,
			},
		})
	}))
	defer server.Close()

	// Create transport
	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	req := Request{
		Method: "POST",
		Path:   "/v1/chat/completions",
		Body: map[string]interface{}{
			"model": "llama3.1-8b",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Hello!",
				},
			},
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = tr.Do(ctx, req)
	}
}

// BenchmarkStream benchmarks SSE streaming with large buffers.
// Expected: 30-40% faster streaming
func BenchmarkStream(b *testing.B) {
	// Create test server that streams responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Send 10 SSE events
		for i := 0; i < 10; i++ {
			data := []byte(`data: {"id":"test","object":"chat.completion.chunk","created":1234567890.123,"model":"llama3.1-8b","choices":[{"index":0,"delta":{"content":"token ` + string(rune('0'+i)) + `"},"finish_reason":null}]}`)
			w.Write(data)
			w.Write([]byte("\n\n"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}

		// Send final event with usage
		data := []byte(`data: {"id":"test","object":"chat.completion.chunk","created":1234567890.123,"model":"llama3.1-8b","choices":[],"usage":{"prompt_tokens":10,"completion_tokens":10,"total_tokens":20}}`)
		w.Write(data)
		w.Write([]byte("\n\n[DONE]\n\n"))
	}))
	defer server.Close()

	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	req := Request{
		Method: "POST",
		Path:   "/v1/chat/completions",
		Body: map[string]interface{}{
			"model": "llama3.1-8b",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Count to 10",
				},
			},
			"stream": true,
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, err := tr.Do(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
		_ = resp
	}
}

// BenchmarkBufferPool benchmarks the buffer pool itself.
// Expected: Near-zero allocations after warmup
func BenchmarkBufferPool(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := getBuffer()
		buf.WriteString("test data")
		_ = buf.Bytes()
		putBuffer(buf)
	}
}

// BenchmarkConnectionReuse benchmarks HTTP connection reuse.
// Expected: >90% connection reuse for repeated requests
func BenchmarkConnectionReuse(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
		})
	}))
	defer server.Close()

	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	req := Request{
		Method: "GET",
		Path:   "/v1/test",
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = tr.Do(ctx, req)
	}
}

// BenchmarkErrorBodyLimit benchmarks error response handling with size limiting.
func BenchmarkErrorBodyLimit(b *testing.B) {
	// Create test server that returns large error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		// Write 100KB error body
		largeError := bytes.Repeat([]byte("x"), 100*1024)
		w.Write(largeError)
	}))
	defer server.Close()

	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	req := Request{
		Method: "POST",
		Path:   "/v1/chat/completions",
		Body:   map[string]interface{}{"model": "test"},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = tr.Do(ctx, req)
	}
}

// Helper functions

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

// Benchmark strings.Builder vs += for string concatenation
func BenchmarkStringConcatenation(b *testing.B) {
	sequences := []string{"seq1", "seq2", "seq3", "seq4"}

	b.Run("strings.Builder", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for j, s := range sequences {
				if j > 0 {
					builder.WriteString(",")
				}
				builder.WriteString(s)
			}
			_ = builder.String()
		}
	})

	b.Run("+= operator", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			joined := ""
			for j, s := range sequences {
				if j > 0 {
					joined += ","
				}
				joined += s
			}
			_ = joined
		}
	})
}
