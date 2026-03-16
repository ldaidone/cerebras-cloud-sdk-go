// Package cerebras provides helper functions for response formats.
package cerebras

// ResponseFormatType represents the type of response format.
type ResponseFormatType string

const (
	// ResponseFormatText specifies plain text output.
	ResponseFormatText ResponseFormatType = "text"
	// ResponseFormatJSONObject specifies JSON object output.
	ResponseFormatJSONObject ResponseFormatType = "json_object"
	// ResponseFormatJSONSchema specifies JSON schema-constrained output.
	ResponseFormatJSONSchema ResponseFormatType = "json_schema"
)

// ResponseFormatJSON creates a ResponseFormat for JSON object output.
//
// Example:
//
//	resp, err := client.ChatCompletions.Create(ctx,
//		Llama31_8b,
//		messages,
//		WithResponseFormat(ResponseFormatJSON()),
//	)
func ResponseFormatJSON() ResponseFormat {
	return ResponseFormat{
		Type: string(ResponseFormatJSONObject),
	}
}

// ResponseFormatJSONWithSchema creates a ResponseFormat for JSON schema output.
//
// Example:
//
//	schema := JSONSchema{
//		Name: "user_info",
//		Description: "Extract user information",
//		Schema: map[string]interface{}{
//			"type": "object",
//			"properties": map[string]interface{}{
//				"name": map[string]interface{}{"type": "string"},
//				"age": map[string]interface{}{"type": "integer"},
//			},
//			"required": []string{"name", "age"},
//		},
//	}
//	resp, err := client.ChatCompletions.Create(ctx,
//		Llama31_8b,
//		messages,
//		WithResponseFormat(ResponseFormatJSONWithSchema(schema)),
//	)
func ResponseFormatJSONWithSchema(schema JSONSchema) ResponseFormat {
	return ResponseFormat{
		Type:       string(ResponseFormatJSONSchema),
		JSONSchema: &schema,
	}
}

// ResponseFormatPlainText creates a ResponseFormat for plain text output.
func ResponseFormatPlainText() ResponseFormat {
	return ResponseFormat{
		Type: string(ResponseFormatText),
	}
}

// DefineJSONSchema creates a JSONSchema with the given name, description, and schema.
//
// Example:
//
//	schema := DefineJSONSchema(
//		"weather_response",
//		"Weather information response",
//		map[string]interface{}{
//			"type": "object",
//			"properties": map[string]interface{}{
//				"location": map[string]interface{}{"type": "string"},
//				"temperature": map[string]interface{}{"type": "number"},
//				"conditions": map[string]interface{}{"type": "string"},
//			},
//			"required": []string{"location", "temperature", "conditions"},
//		},
//	)
func DefineJSONSchema(name, description string, schema interface{}) JSONSchema {
	return JSONSchema{
		Name:        name,
		Description: description,
		Schema:      schema,
	}
}

// JSONSchemaBuilder provides a fluent interface for building JSON schemas.
type JSONSchemaBuilder struct {
	schema map[string]interface{}
}

// NewJSONSchemaBuilder creates a new JSON schema builder.
func NewJSONSchemaBuilder() *JSONSchemaBuilder {
	return &JSONSchemaBuilder{
		schema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}
}

// WithProperty adds a property to the schema.
func (b *JSONSchemaBuilder) WithProperty(name, propType string, opts ...SchemaOption) *JSONSchemaBuilder {
	props := b.schema["properties"].(map[string]interface{})
	props[name] = JSONSchemaProperty(propType, opts...)
	return b
}

// WithRequired adds required properties to the schema.
func (b *JSONSchemaBuilder) WithRequired(props ...string) *JSONSchemaBuilder {
	b.schema["required"] = props
	return b
}

// WithDescription adds a description to the schema.
func (b *JSONSchemaBuilder) WithDescription(desc string) *JSONSchemaBuilder {
	b.schema["description"] = desc
	return b
}

// Build returns the built schema.
func (b *JSONSchemaBuilder) Build() map[string]interface{} {
	return b.schema
}

// BuildWithSchema returns a JSONSchema with the built schema.
func (b *JSONSchemaBuilder) BuildWithSchema(name, description string) JSONSchema {
	return JSONSchema{
		Name:        name,
		Description: description,
		Schema:      b.schema,
	}
}
