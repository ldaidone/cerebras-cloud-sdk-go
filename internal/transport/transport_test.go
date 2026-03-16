package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTransport(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		tr := NewTransport(TransportConfig{})

		if tr.maxRetries != DefaultMaxRetries {
			t.Errorf("maxRetries = %d, want %d", tr.maxRetries, DefaultMaxRetries)
		}
		if tr.timeout != DefaultTimeout {
			t.Errorf("timeout = %v, want %v", tr.timeout, DefaultTimeout)
		}
		if tr.httpClient == nil {
			t.Error("httpClient should not be nil")
		}
	})

	t.Run("custom configuration", func(t *testing.T) {
		customClient := &http.Client{Timeout: 30 * time.Second}
		tr := NewTransport(TransportConfig{
			HTTPClient: customClient,
			APIKey:     "test-key",
			BaseURL:    "https://test.example.com",
			MaxRetries: 5,
			Timeout:    30 * time.Second,
		})

		if tr.maxRetries != 5 {
			t.Errorf("maxRetries = %d, want 5", tr.maxRetries)
		}
		if tr.timeout != 30*time.Second {
			t.Errorf("timeout = %v, want 30s", tr.timeout)
		}
		if tr.httpClient != customClient {
			t.Error("httpClient should be the custom client")
		}
		if tr.apiKey != "test-key" {
			t.Errorf("apiKey = %s, want test-key", tr.apiKey)
		}
		if tr.baseURL != "https://test.example.com" {
			t.Errorf("baseURL = %s, want https://test.example.com", tr.baseURL)
		}
	})
}

func TestTransport_Do_Success(t *testing.T) {
	type ResponseBody struct {
		Message string `json:"message"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/v1/test" {
			t.Errorf("Path = %s, want /v1/test", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("Authorization = %s, want Bearer test-key", auth)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", ct)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseBody{Message: "success"})
	}))
	defer server.Close()

	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	req := Request{
		Method: http.MethodPost,
		Path:   "/v1/test",
		Body:   map[string]string{"test": "data"},
	}

	resp, err := tr.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body ResponseBody
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if body.Message != "success" {
		t.Errorf("Message = %s, want success", body.Message)
	}
}

func TestTransport_Do_ErrorStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantRetry  bool
	}{
		{"400 Bad Request", http.StatusBadRequest, `{"error": "bad request"}`, false},
		{"401 Unauthorized", http.StatusUnauthorized, `{"error": "unauthorized"}`, false},
		{"403 Forbidden", http.StatusForbidden, `{"error": "forbidden"}`, false},
		{"404 Not Found", http.StatusNotFound, `{"error": "not found"}`, false},
		{"429 Rate Limit", http.StatusTooManyRequests, `{"error": "rate limited"}`, true},
		{"500 Internal Error", http.StatusInternalServerError, `{"error": "server error"}`, true},
		{"503 Unavailable", http.StatusServiceUnavailable, `{"error": "unavailable"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			tr := NewTransport(TransportConfig{
				APIKey:     "test-key",
				BaseURL:    server.URL,
				MaxRetries: 0, // Disable retries for this test
			})

			req := Request{
				Method: http.MethodGet,
				Path:   "/v1/test",
			}

			_, err := tr.Do(context.Background(), req)
			if err == nil {
				t.Fatal("Do() should return an error")
			}
		})
	}
}

func TestTransport_Do_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := Request{
		Method: http.MethodGet,
		Path:   "/v1/test",
	}

	_, err := tr.Do(ctx, req)
	if err == nil {
		t.Fatal("Do() should return an error for cancelled context")
	}
}

func TestTransport_DoJSON(t *testing.T) {
	type ResponseBody struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseBody{ID: 1, Name: "test"})
	}))
	defer server.Close()

	tr := NewTransport(TransportConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	req := Request{
		Method: http.MethodGet,
		Path:   "/v1/test",
	}

	var result ResponseBody
	err := tr.DoJSON(context.Background(), req, &result)
	if err != nil {
		t.Fatalf("DoJSON() error = %v", err)
	}

	if result.ID != 1 {
		t.Errorf("ID = %d, want 1", result.ID)
	}
	if result.Name != "test" {
		t.Errorf("Name = %s, want test", result.Name)
	}
}

func TestTransport_CalculateBackoff(t *testing.T) {
	tr := NewTransport(TransportConfig{})

	t.Run("exponential backoff increases", func(t *testing.T) {
		backoff0 := tr.calculateBackoff(0, nil)
		backoff1 := tr.calculateBackoff(1, nil)
		backoff2 := tr.calculateBackoff(2, nil)

		if backoff1 <= backoff0 {
			t.Error("backoff should increase with attempt number")
		}
		if backoff2 <= backoff1 {
			t.Error("backoff should increase with attempt number")
		}
	})

	t.Run("max backoff ceiling", func(t *testing.T) {
		backoff := tr.calculateBackoff(10, nil)
		if backoff > DefaultMaxBackoff {
			t.Errorf("backoff = %v, should be capped at %v", backoff, DefaultMaxBackoff)
		}
	})

	t.Run("Retry-After header", func(t *testing.T) {
		resp := &Response{
			Header: http.Header{"Retry-After": []string{"30"}},
		}
		backoff := tr.calculateBackoff(0, resp)
		
		expected := 30 * time.Second
		if backoff != expected {
			t.Errorf("backoff with Retry-After = %v, want %v", backoff, expected)
		}
	})
}

func TestTransport_Getters(t *testing.T) {
	customClient := &http.Client{Timeout: 45 * time.Second}
	tr := NewTransport(TransportConfig{
		HTTPClient: customClient,
		APIKey:     "test-key",
		BaseURL:    "https://test.example.com",
		Timeout:    45 * time.Second,
	})

	if got := tr.GetBaseURL(); got != "https://test.example.com" {
		t.Errorf("GetBaseURL() = %s, want https://test.example.com", got)
	}
	if got := tr.GetTimeout(); got != 45*time.Second {
		t.Errorf("GetTimeout() = %v, want 45s", got)
	}
	if got := tr.GetHTTPClient(); got != customClient {
		t.Error("GetHTTPClient() should return the custom client")
	}
}
