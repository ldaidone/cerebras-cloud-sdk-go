package cerebras

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestModelsService_List(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/models" {
			t.Errorf("Path = %s, want /v1/models", r.URL.Path)
		}

		// Verify headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("Authorization = %s, want Bearer test-key", auth)
		}

		// Send response
		response := ModelList{
			Object: "list",
			Data: []Model{
				{
					ID:      Llama31_8b,
					Object:  "model",
					Created: time.Now().Unix(),
					OwnedBy: "meta",
				},
				{
					ID:      Llama31_70b,
					Object:  "model",
					Created: time.Now().Unix(),
					OwnedBy: "meta",
				},
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

	// Test List method
	ctx := context.Background()
	models, err := client.Models.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if models.Object != "list" {
		t.Errorf("Object = %s, want list", models.Object)
	}
	if len(models.Data) != 2 {
		t.Errorf("Data length = %d, want 2", len(models.Data))
	}
	if models.Data[0].ID != Llama31_8b {
		t.Errorf("First model ID = %s, want %s", models.Data[0].ID, Llama31_8b)
	}
	if models.Data[0].OwnedBy != "meta" {
		t.Errorf("First model OwnedBy = %s, want meta", models.Data[0].OwnedBy)
	}
}

func TestModelsService_Retrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/models/llama3.1-8b" {
			t.Errorf("Path = %s, want /v1/models/llama3.1-8b", r.URL.Path)
		}

		// Send response
		response := Model{
			ID:      Llama31_8b,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "meta",
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
	model, err := client.Models.Retrieve(ctx, Llama31_8b)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}

	if model.ID != Llama31_8b {
		t.Errorf("ID = %s, want %s", model.ID, Llama31_8b)
	}
	if model.OwnedBy != "meta" {
		t.Errorf("OwnedBy = %s, want meta", model.OwnedBy)
	}
}

func TestModelsService_Retrieve_Validation(t *testing.T) {
	client := NewClient(
		WithAPIKey("test-key"),
		WithTCPWarming(false),
	)

	ctx := context.Background()

	t.Run("missing modelID", func(t *testing.T) {
		_, err := client.Models.Retrieve(ctx, "")
		if err == nil {
			t.Error("Retrieve() should return error for missing modelID")
		}
	})
}

func TestModelsService_Retrieve_404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": {"message": "Model not found"}}`))
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithTCPWarming(false),
		WithMaxRetries(0), // Disable retries for this test
	)

	ctx := context.Background()
	_, err := client.Models.Retrieve(ctx, "nonexistent-model")
	if err == nil {
		t.Fatal("Retrieve() should return error for 404")
	}
}

func TestModelsService_List_ContextCancellation(t *testing.T) {
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

	_, err := client.Models.List(ctx)
	if err == nil {
		t.Fatal("List() should return error for cancelled context")
	}
}

func TestModelsService_Retrieve_ContextCancellation(t *testing.T) {
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

	_, err := client.Models.Retrieve(ctx, Llama31_8b)
	if err == nil {
		t.Fatal("Retrieve() should return error for cancelled context")
	}
}

func TestModelsService_ErrorHandling(t *testing.T) {
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

		_, err := client.Models.List(context.Background())
		if err == nil {
			t.Fatal("List() should return error for 401")
		}
	})

	t.Run("500 Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": {"message": "Server error"}}`))
		}))
		defer server.Close()

		client := NewClient(
			WithAPIKey("test-key"),
			WithBaseURL(server.URL),
			WithTCPWarming(false),
			WithMaxRetries(0), // Disable retries for this test
		)

		_, err := client.Models.List(context.Background())
		if err == nil {
			t.Fatal("List() should return error for 500")
		}
	})
}
