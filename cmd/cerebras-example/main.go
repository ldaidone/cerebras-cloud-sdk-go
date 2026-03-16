// Cerebras Cloud SDK Example
//
// This example demonstrates how to use the Cerebras Cloud SDK for Go.
// It reads the API key from the CEREBRAS_API_KEY environment variable.
//
// Usage:
//
//	1. Set your API key:
//	   export CEREBRAS_API_KEY=sk-your-api-key-here
//
//	2. Run the example:
//	   go run cmd/cerebras-example/main.go
//
//	3. Try different examples by passing a flag:
//	   go run cmd/cerebras-example/main.go -example=models
//	   go run cmd/cerebras-example/main.go -example=chat
//	   go run cmd/cerebras-example/main.go -example=stream
//	   go run cmd/cerebras-example/main.go -example=text
//	   go run cmd/cerebras-example/main.go -example=tools
//	   go run cmd/cerebras-example/main.go -example=json
//
// Available examples:
//   - models: List and retrieve models
//   - chat: Basic chat completion
//   - stream: Streaming chat completion
//   - text: Text completion (legacy)
//   - tools: Tool/function calling
//   - json: JSON response format
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github/ldaidone/cerebras-cloud-sdk-go/pkg/cerebras"
)

func main() {
	// Parse command-line flags
	example := flag.String("example", "chat", "Example to run: models, chat, stream, text, tools, json")
	timeout := flag.Duration("timeout", 30*time.Second, "Request timeout")
	flag.Parse()

	// Check if API key is set
	if os.Getenv("CEREBRAS_API_KEY") == "" {
		log.Fatal("Error: CEREBRAS_API_KEY environment variable is not set\n" +
			"Please set it with: export CEREBRAS_API_KEY=sk-your-api-key-here")
	}

	// Create client with environment variable API key
	// (API key is automatically read from CEREBRAS_API_KEY)
	client := cerebras.NewClient(
		cerebras.WithTimeout(*timeout),
		// Optional: Set custom base URL for local development
		// cerebras.WithBaseURL("http://localhost:11434"),
		// cerebras.WithTCPWarming(false),
	)

	fmt.Printf("Running example: %s\n\n", *example)

	// Run the selected example
	var err error
	switch *example {
	case "models":
		err = runModelsExample(client)
	case "chat":
		err = runChatExample(client)
	case "stream":
		err = runStreamExample(client)
	case "text":
		err = runTextExample(client)
	case "tools":
		err = runToolsExample(client)
	case "json":
		err = runJSONExample(client)
	default:
		log.Fatalf("Unknown example: %s\nAvailable examples: models, chat, stream, text, tools, json", *example)
	}

	if err != nil {
		log.Fatalf("Example failed: %v", err)
	}

	fmt.Println("\n✅ Example completed successfully!")
}

// runModelsExample demonstrates listing and retrieving models
func runModelsExample(client *cerebras.Client) error {
	fmt.Println("=== Models Example ===\n")

	// List all available models
	fmt.Println("1. Listing all models...")
	models, err := client.Models.List(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	fmt.Printf("   Found %d models:\n", len(models.Data))
	for _, model := range models.Data {
		fmt.Printf("   - %s (owned by %s)\n", model.ID, model.OwnedBy)
	}
	fmt.Println()

	// Retrieve a specific model
	fmt.Println("2. Retrieving model details...")
	model, err := client.Models.Retrieve(context.Background(), cerebras.Llama31_8b)
	if err != nil {
		return fmt.Errorf("failed to retrieve model: %w", err)
	}

	fmt.Printf("   Model: %s\n", model.ID)
	fmt.Printf("   Owned by: %s\n", model.OwnedBy)
	fmt.Printf("   Created: %d\n", model.Created)
	fmt.Println()

	return nil
}

// runChatExample demonstrates basic chat completion
func runChatExample(client *cerebras.Client) error {
	fmt.Println("=== Chat Completion Example ===\n")

	// Create a chat completion
	fmt.Println("Sending chat completion request...")
	response, err := client.ChatCompletions.Create(
		context.Background(),
		cerebras.Llama31_8b,
		[]cerebras.Message{
			{Role: cerebras.MessageRoleSystem, Content: "You are a helpful assistant."},
			{Role: cerebras.MessageRoleUser, Content: "What is the capital of France?"},
		},
		cerebras.WithTemperature(0.7),
		cerebras.WithMaxTokens(100),
	)
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}

	// Display response
	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Usage: %d prompt tokens, %d completion tokens, %d total\n",
		response.Usage.PromptTokens,
		response.Usage.CompletionTokens,
		response.Usage.TotalTokens,
	)
	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Println()

	return nil
}

// runStreamExample demonstrates streaming chat completion
func runStreamExample(client *cerebras.Client) error {
	fmt.Println("=== Streaming Chat Completion Example ===\n")

	fmt.Println("Sending streaming chat completion request...")
	fmt.Println("Response: ")
	fmt.Print("> ")

	// Create streaming chat completion
	stream, errChan, err := client.ChatCompletions.CreateStream(
		context.Background(),
		cerebras.Llama31_8b,
		[]cerebras.Message{
			{Role: cerebras.MessageRoleSystem, Content: "You are a helpful assistant."},
			{Role: cerebras.MessageRoleUser, Content: "Count from 1 to 5."},
		},
		cerebras.WithTemperature(0.7),
		cerebras.WithMaxTokens(50),
	)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	// Read stream
	for {
		select {
		case chunk, ok := <-stream:
			if !ok {
				// Stream closed
				fmt.Println("\n")
				return nil
			}
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
				fmt.Print(*chunk.Choices[0].Delta.Content)
			}
		case err := <-errChan:
			return fmt.Errorf("stream error: %w", err)
		}
	}
}

// runTextExample demonstrates text completion (legacy endpoint)
func runTextExample(client *cerebras.Client) error {
	fmt.Println("=== Text Completion Example (Legacy) ===\n")

	fmt.Println("Sending text completion request...")
	response, err := client.TextCompletions.Create(
		context.Background(),
		cerebras.Llama31_8b,
		"Once upon a time in a land far away,",
		cerebras.WithMaxTokensText(50),
		cerebras.WithTemperatureText(0.7),
	)
	if err != nil {
		return fmt.Errorf("failed to create text completion: %w", err)
	}

	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Usage: %d prompt tokens, %d completion tokens\n",
		response.Usage.PromptTokens,
		response.Usage.CompletionTokens,
	)
	fmt.Printf("Completion: %s\n", response.Choices[0].Text)
	fmt.Println()

	return nil
}

// runToolsExample demonstrates tool/function calling
func runToolsExample(client *cerebras.Client) error {
	fmt.Println("=== Tool Calling Example ===\n")

	// Define a tool
	weatherTool := cerebras.DefineTool(
		"get_weather",
		"Get the current weather in a given location",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city and state, e.g. San Francisco, CA",
				},
				"unit": map[string]interface{}{
					"type": "string",
					"enum": []string{"celsius", "fahrenheit"},
				},
			},
			"required": []string{"location"},
		},
	)

	fmt.Println("Sending chat completion with tool...")
	response, err := client.ChatCompletions.Create(
		context.Background(),
		cerebras.Llama31_8b,
		[]cerebras.Message{
			{Role: cerebras.MessageRoleUser, Content: "What's the weather in San Francisco?"},
		},
		cerebras.WithTools(weatherTool),
		cerebras.WithToolChoice(cerebras.ToolChoiceAuto()),
	)
	if err != nil {
		return fmt.Errorf("failed to create chat with tools: %w", err)
	}

	// Check if model called a tool
	if len(response.Choices) > 0 && len(response.Choices[0].Message.ToolCalls) > 0 {
		toolCall := response.Choices[0].Message.ToolCalls[0]
		fmt.Printf("Model called tool: %s\n", toolCall.Function.Name)
		fmt.Printf("Tool arguments: %s\n", toolCall.Function.Arguments)

		// Parse tool arguments
		var args struct {
			Location string `json:"location"`
			Unit     string `json:"unit"`
		}
		if err := cerebras.ParseFunctionArguments(toolCall, &args); err != nil {
			return fmt.Errorf("failed to parse arguments: %w", err)
		}
		fmt.Printf("Parsed location: %s, unit: %s\n", args.Location, args.Unit)
	} else {
		fmt.Printf("Model response: %s\n", response.Choices[0].Message.Content)
	}
	fmt.Println()

	return nil
}

// runJSONExample demonstrates JSON response format
func runJSONExample(client *cerebras.Client) error {
	fmt.Println("=== JSON Response Format Example ===\n")

	// Define JSON schema
	schema := cerebras.DefineJSONSchema(
		"user_info",
		"Extract user information from text",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "The person's name",
				},
				"age": map[string]interface{}{
					"type":        "integer",
					"description": "The person's age",
				},
				"city": map[string]interface{}{
					"type":        "string",
					"description": "The person's city",
				},
			},
			"required": []string{"name", "age", "city"},
		},
	)

	fmt.Println("Sending chat completion with JSON schema...")
	response, err := client.ChatCompletions.Create(
		context.Background(),
		cerebras.Llama31_8b,
		[]cerebras.Message{
			{Role: cerebras.MessageRoleSystem, Content: "Extract information from the text."},
			{Role: cerebras.MessageRoleUser, Content: "John is 30 years old and lives in New York."},
		},
		cerebras.WithResponseFormat(cerebras.ResponseFormatJSONWithSchema(schema)),
		cerebras.WithMaxTokens(100),
	)
	if err != nil {
		return fmt.Errorf("failed to create JSON chat: %w", err)
	}

	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("JSON Response: %s\n", response.Choices[0].Message.Content)
	fmt.Println()

	return nil
}
