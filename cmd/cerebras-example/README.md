# Cerebras Cloud SDK Example

This example demonstrates how to use the Cerebras Cloud SDK for Go with various features.

## Prerequisites

1. **Set your API key:**
   ```bash
   export CEREBRAS_API_KEY=sk-your-api-key-here
   ```

2. **Build the example:**
   ```bash
   cd /Volumes/Mac\ HD\ 1TB/Users/leodaido/projects/LEO/cerebras-cloud-sdk-go
   go build -o cerebras-example ./cmd/cerebras-example/
   ```

## Usage

Run the example with different features:

### 1. Models Example
List and retrieve available models:
```bash
./cerebras-example -example=models
```

### 2. Chat Completion Example
Basic non-streaming chat completion:
```bash
./cerebras-example -example=chat
```

### 3. Streaming Example
Streaming chat completion with real-time tokens:
```bash
./cerebras-example -example=stream
```

### 4. Text Completion Example
Legacy text completion endpoint:
```bash
./cerebras-example -example=text
```

### 5. Tool Calling Example
Function/tool calling with JSON schema:
```bash
./cerebras-example -example=tools
```

### 6. JSON Response Example
Structured JSON output with schema:
```bash
./cerebras-example -example=json
```

## Command-Line Options

```bash
./cerebras-example [options]

Options:
  -example string
        Example to run: models, chat, stream, text, tools, json (default "chat")
  -timeout duration
        Request timeout (default 30s)
```

## Examples

### Basic Chat
```bash
./cerebras-example -example=chat -timeout=60s
```

### Streaming with Custom Timeout
```bash
./cerebras-example -example=stream -timeout=120s
```

## Running with `go run`

You can also run directly without building:

```bash
go run cmd/cerebras-example/main.go -example=chat
go run cmd/cerebras-example/main.go -example=stream
go run cmd/cerebras-example/main.go -example=models
```

## What Each Example Demonstrates

### Models (`-example=models`)
- Listing all available models
- Retrieving specific model details
- Using model constants (Llama31_8b, etc.)

### Chat (`-example=chat`)
- Basic chat completion
- System and user messages
- Temperature and max_tokens parameters
- Usage statistics

### Stream (`-example=stream`)
- Channel-based streaming
- Real-time token delivery
- Error handling during streaming
- Context cancellation

### Text (`-example=text`)
- Legacy completions endpoint
- Prompt-based completion
- Text-specific parameters

### Tools (`-example=tools`)
- Tool/function definition
- JSON Schema for parameters
- Tool choice (auto, none, required)
- Parsing tool call arguments

### JSON (`-example=json`)
- JSON object response format
- JSON schema-constrained output
- Structured data extraction

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CEREBRAS_API_KEY` | Your Cerebras Cloud API key | ✅ Yes |
| `CEREBRAS_BASE_URL` | Custom API base URL (optional) | ❌ No |

## Error Handling

The example demonstrates proper error handling:
- API key validation
- Request timeouts
- Network errors
- API error responses

## Next Steps

After trying the examples, you can:
1. Copy the code patterns into your project
2. Modify the examples for your use case
3. Explore the full API in `pkg/cerebras/`
4. Read the full documentation at https://pkg.go.dev/github.com/ldaidone/cerebras-cloud-sdk-go

## Support

For issues or questions:
- Check the main README.md
- Review GoDoc: https://pkg.go.dev/github.com/ldaidone/cerebras-cloud-sdk-go/pkg/cerebras
- See the COMPLETE.md for full feature list
