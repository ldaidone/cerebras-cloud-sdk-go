package cerebras_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github/ldaidone/cerebras-cloud-sdk-go/pkg/cerebras"
)

// ExampleNewClient demonstrates basic client creation.
func ExampleNewClient() {
	// Basic client with API key
	_ = cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
	)

	fmt.Println("Client created")
}

// ExampleNewClient_withOptions demonstrates client creation with multiple options.
func ExampleNewClient_withOptions() {
	client := cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
		cerebras.WithBaseURL("https://api.cerebras.ai"),
		cerebras.WithTimeout(30*time.Second),
		cerebras.WithMaxRetries(5),
		cerebras.WithTCPWarming(true),
	)

	fmt.Println("Client configured with timeout:", client.Timeout())
}

// ExampleNewClient_withEnvironment demonstrates using environment variables.
// Set CEREBRAS_API_KEY and CEREBRAS_BASE_URL environment variables before running.
func ExampleNewClient_withEnvironment() {
	// Client will use CEREBRAS_API_KEY and CEREBRAS_BASE_URL from environment
	client := cerebras.NewClient()

	fmt.Println("Client created with environment configuration")
	fmt.Println("Base URL:", client.BaseURL())
}

// ExampleNewClient_localDevelopment demonstrates configuration for local development.
func ExampleNewClient_localDevelopment() {
	// For local development with Ollama, LM Studio, or other proxies
	_ = cerebras.NewClient(
		cerebras.WithBaseURL("http://localhost:11434"),
		cerebras.WithTCPWarming(false), // Disable warming for local endpoints
		cerebras.WithTimeout(10*time.Second),
	)

	fmt.Println("Local development client created")
}

// ExampleNewClient_customHTTPClient demonstrates using a custom HTTP client.
func ExampleNewClient_customHTTPClient() {
	customClient := &http.Client{
		Timeout: 45 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	client := cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
		cerebras.WithHTTPClient(customClient),
	)

	fmt.Println("Client created with custom HTTP client")
	_ = client
}

// ExampleClient_Get demonstrates making a GET request.
func ExampleClient_Get() {
	client := cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result map[string]interface{}
	err := client.Get(ctx, "/v1/models", &result)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Println("Request succeeded")
}

// ExampleClient_Post demonstrates making a POST request.
func ExampleClient_Post() {
	client := cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	requestBody := map[string]interface{}{
		"model": "llama3.1-8b-8192",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello, world!"},
		},
	}

	var result map[string]interface{}
	err := client.Post(ctx, "/v1/chat/completions", requestBody, &result)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Println("Chat completion created")
}

// ExampleClient_Do demonstrates making a custom HTTP request.
func ExampleClient_Do() {
	client := cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result map[string]interface{}
	err := client.Do(ctx, http.MethodDelete, "/v1/some-resource", nil, &result)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}

	fmt.Println("Custom request succeeded")
}

// ExampleClient_contextCancellation demonstrates context cancellation.
func ExampleClient_contextCancellation() {
	client := cerebras.NewClient(
		cerebras.WithAPIKey("sk-your-api-key-here"),
	)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Simulate some work
	go func() {
		time.Sleep(5 * time.Second)
		cancel() // Cancel after 5 seconds
	}()

	var result map[string]interface{}
	err := client.Get(ctx, "/v1/slow-endpoint", &result)
	if err != nil {
		log.Printf("Request cancelled or failed: %v", err)
		return
	}

	fmt.Println("Request completed before cancellation")
}
