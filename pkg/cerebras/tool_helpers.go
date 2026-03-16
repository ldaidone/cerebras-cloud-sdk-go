// Package cerebras provides helper functions for tool calling.
package cerebras

import "encoding/json"

// DefineFunction creates a Function definition with the given name, description, and parameters.
// This is a convenience function for defining tools.
//
// Example:
//
//	weatherFunc := DefineFunction(
//		"get_weather",
//		"Get the current weather in a given location",
//		map[string]interface{}{
//			"type": "object",
//			"properties": map[string]interface{}{
//				"location": map[string]interface{}{
//					"type": "string",
//					"description": "The city and state, e.g. San Francisco, CA",
//				},
//				"unit": map[string]interface{}{
//					"type": "string",
//					"enum": []string{"celsius", "fahrenheit"},
//				},
//			},
//			"required": []string{"location"},
//		},
//	)
func DefineFunction(name, description string, parameters interface{}) Function {
	return Function{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}
}

// DefineTool creates a Tool with the given function definition.
// This is a convenience function that wraps DefineFunction.
//
// Example:
//
//	weatherTool := DefineTool(
//		"get_weather",
//		"Get the current weather in a given location",
//		map[string]interface{}{
//			"type": "object",
//			"properties": map[string]interface{}{
//				"location": map[string]interface{}{
//					"type": "string",
//					"description": "The city and state",
//				},
//			},
//			"required": []string{"location"},
//		},
//	)
func DefineTool(name, description string, parameters interface{}) Tool {
	return Tool{
		Type:     "function",
		Function: DefineFunction(name, description, parameters),
	}
}

// JSONSchemaProperty creates a JSON Schema property definition.
//
// Example:
//
//	locationProp := JSONSchemaProperty("string",
//		WithDescription("The city and state, e.g. San Francisco, CA"),
//		WithRequired(true),
//	)
func JSONSchemaProperty(propType string, opts ...SchemaOption) map[string]interface{} {
	prop := map[string]interface{}{
		"type": propType,
	}

	for _, opt := range opts {
		opt(prop)
	}

	return prop
}

// SchemaOption is a functional option for JSON Schema properties.
type SchemaOption func(map[string]interface{})

// WithDescription adds a description to a JSON Schema property.
func WithDescription(desc string) SchemaOption {
	return func(prop map[string]interface{}) {
		prop["description"] = desc
	}
}

// WithRequired marks a JSON Schema property as required.
// Note: This adds the property to the parent's "required" array.
func WithRequired(required bool) SchemaOption {
	return func(prop map[string]interface{}) {
		if required {
			prop["required"] = true
		}
	}
}

// WithEnum adds enum values to a JSON Schema property.
func WithEnum(values ...interface{}) SchemaOption {
	return func(prop map[string]interface{}) {
		if len(values) > 0 {
			prop["enum"] = values
		}
	}
}

// WithDefault sets a default value for a JSON Schema property.
func WithDefault(value interface{}) SchemaOption {
	return func(prop map[string]interface{}) {
		prop["default"] = value
	}
}

// JSONSchemaObject creates a JSON Schema object with properties.
//
// Example:
//
//	schema := JSONSchemaObject(
//		JSONSchemaProperty("string",
//			WithDescription("The city and state"),
//		),
//	)
func JSONSchemaObject(properties map[string]interface{}, required ...string) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// ToolChoiceAuto returns the "auto" tool choice.
// The model will decide whether to use tools.
func ToolChoiceAuto() string {
	return "auto"
}

// ToolChoiceNone returns the "none" tool choice.
// The model will not use any tools.
func ToolChoiceNone() string {
	return "none"
}

// ToolChoiceRequired returns the "required" tool choice.
// The model must use one or more tools.
func ToolChoiceRequired() string {
	return "required"
}

// ToolChoiceFunction returns a tool choice for a specific function.
//
// Example:
//
//	choice := ToolChoiceFunc("get_weather")
func ToolChoiceFunc(name string) ToolChoice {
	return ToolChoice{
		Type: "function",
		Function: ToolChoiceFunction{
			Name: name,
		},
	}
}

// ParseFunctionArguments parses function arguments from a tool call.
// This is useful for extracting typed parameters from the raw JSON arguments.
//
// Example:
//
//	var args struct {
//		Location string `json:"location"`
//		Unit     string `json:"unit"`
//	}
//	err := ParseFunctionArguments(toolCall, &args)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Getting weather for %s\n", args.Location)
func ParseFunctionArguments(toolCall ToolCall, target interface{}) error {
	return json.Unmarshal([]byte(toolCall.Function.Arguments), target)
}
