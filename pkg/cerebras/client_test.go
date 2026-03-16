package cerebras

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	// Clear env vars to ensure clean test
	os.Unsetenv("CEREBRAS_API_KEY")
	os.Unsetenv("CEREBRAS_BASE_URL")

	client := NewClient()

	if client.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %s, want %s", client.baseURL, DefaultBaseURL)
	}
	if client.timeout != DefaultTimeout {
		t.Errorf("timeout = %v, want %v", client.timeout, DefaultTimeout)
	}
	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("maxRetries = %d, want %d", client.maxRetries, DefaultMaxRetries)
	}
	if client.tcpWarming != true {
		t.Error("tcpWarming should be true by default")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 45 * time.Second}

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL("https://test.example.com"),
		WithTimeout(30*time.Second),
		WithMaxRetries(5),
		WithTCPWarming(false),
		WithHTTPClient(customClient),
	)

	if client.apiKey != "test-key" {
		t.Errorf("apiKey = %s, want test-key", client.apiKey)
	}
	if client.baseURL != "https://test.example.com" {
		t.Errorf("baseURL = %s, want https://test.example.com", client.baseURL)
	}
	if client.timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", client.timeout)
	}
	if client.maxRetries != 5 {
		t.Errorf("maxRetries = %d, want 5", client.maxRetries)
	}
	if client.tcpWarming != false {
		t.Error("tcpWarming should be false")
	}
	if client.httpClient != customClient {
		t.Error("httpClient should be the custom client")
	}
}

func TestNewClient_EnvironmentVariables(t *testing.T) {
	// Set env vars
	os.Setenv("CEREBRAS_API_KEY", "env-key")
	os.Setenv("CEREBRAS_BASE_URL", "https://env.example.com")
	defer func() {
		os.Unsetenv("CEREBRAS_API_KEY")
		os.Unsetenv("CEREBRAS_BASE_URL")
	}()

	client := NewClient()

	if client.apiKey != "env-key" {
		t.Errorf("apiKey = %s, want env-key", client.apiKey)
	}
	if client.baseURL != "https://env.example.com" {
		t.Errorf("baseURL = %s, want https://env.example.com", client.baseURL)
	}
}

func TestNewClient_ExplicitOverridesEnv(t *testing.T) {
	// Set env vars
	os.Setenv("CEREBRAS_API_KEY", "env-key")
	os.Setenv("CEREBRAS_BASE_URL", "https://env.example.com")
	defer func() {
		os.Unsetenv("CEREBRAS_API_KEY")
		os.Unsetenv("CEREBRAS_BASE_URL")
	}()

	client := NewClient(
		WithAPIKey("explicit-key"),
		WithBaseURL("https://explicit.example.com"),
	)

	if client.apiKey != "explicit-key" {
		t.Errorf("apiKey = %s, want explicit-key", client.apiKey)
	}
	if client.baseURL != "https://explicit.example.com" {
		t.Errorf("baseURL = %s, want https://explicit.example.com", client.baseURL)
	}
}

func TestClient_Accessors(t *testing.T) {
	client := NewClient(
		WithBaseURL("https://test.example.com"),
		WithTimeout(45 * time.Second),
	)

	if got := client.BaseURL(); got != "https://test.example.com" {
		t.Errorf("BaseURL() = %s, want https://test.example.com", got)
	}
	if got := client.Timeout(); got != 45*time.Second {
		t.Errorf("Timeout() = %v, want 45s", got)
	}
}

func TestClient_Do_Methods(t *testing.T) {
	// Create a mock transport
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithTCPWarming(false),
	)

	t.Run("Get", func(t *testing.T) {
		var result map[string]interface{}
		err := client.Get(context.Background(), "/v1/test", &result)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if result["status"] != "ok" {
			t.Errorf("status = %v, want ok", result["status"])
		}
	})

	t.Run("Post", func(t *testing.T) {
		body := map[string]string{"data": "test"}
		var result map[string]interface{}
		err := client.Post(context.Background(), "/v1/test", body, &result)
		if err != nil {
			t.Fatalf("Post() error = %v", err)
		}
		if result["status"] != "ok" {
			t.Errorf("status = %v, want ok", result["status"])
		}
	})
}

func TestClient_Do_ContextCancellation(t *testing.T) {
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

	var result map[string]interface{}
	err := client.Get(ctx, "/v1/test", &result)
	if err == nil {
		t.Fatal("Get() should return an error for cancelled context")
	}
}
