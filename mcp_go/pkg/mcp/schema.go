package mcp

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// SchemaToOptions Convert struct to MCP ToolOption list
func SchemaToOptions(schema any) ([]mcp.ToolOption, error) {
	t := reflect.TypeOf(schema)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var options []mcp.ToolOption
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		desc := field.Tag.Get("desc")
		required := field.Tag.Get("required") == "true"
		enumTag := field.Tag.Get("enum")
		opts := []mcp.PropertyOption{}

		if desc != "" {
			opts = append(opts, mcp.Description(desc))
		}
		if required {
			opts = append(opts, mcp.Required())
		}
		if enumTag != "" && field.Type.Kind() == reflect.String {
			enumValues := strings.Split(enumTag, ",")
			opts = append(opts, mcp.Enum(enumValues...))
		}

		switch field.Type.Kind() {
		case reflect.Int:
			options = append(options, mcp.WithNumber(jsonTag, opts...))
		case reflect.Bool:
			options = append(options, mcp.WithBoolean(jsonTag, opts...))
		case reflect.String:
			options = append(options, mcp.WithString(jsonTag, opts...))
		case reflect.Struct:
			subSchema := reflect.New(field.Type).Interface()
			subOptions, err := SchemaToOptions(subSchema)
			if err != nil {
				return nil, err
			}

			// Create a temporary Tool to get JSON Schema of sub-struct
			tempTool := mcp.NewTool("temp", subOptions...)
			tempJSON, _ := tempTool.MarshalJSON()
			var tempMap map[string]any
			if err := json.Unmarshal(tempJSON, &tempMap); err != nil {
				continue
			}

			// Extract properties from temporary Tool
			if inputSchema, ok := tempMap["inputSchema"].(map[string]any); ok {
				if properties, ok := inputSchema["properties"].(map[string]any); ok {
					// Check if there are required fields
					if required, ok := inputSchema["required"].([]any); ok {
						// Add required field information to corresponding properties
						for _, req := range required {
							if reqStr, ok := req.(string); ok {
								if prop, ok := properties[reqStr].(map[string]any); ok {
									prop["required"] = true
								}
							}
						}
					}
					opts = append(opts, mcp.Properties(properties))
				}
			}
			options = append(options, mcp.WithObject(jsonTag, opts...))

		case reflect.Slice:
			elemType := field.Type.Elem()
			var items map[string]any

			switch elemType.Kind() {
			case reflect.Int:
				items = map[string]any{
					"type": "number",
				}
			case reflect.Bool:
				items = map[string]any{
					"type": "boolean",
				}
			case reflect.String:
				items = map[string]any{
					"type": "string",
				}
			case reflect.Struct:
				subSchema := reflect.New(elemType).Interface()
				subOptions, err := SchemaToOptions(subSchema)
				if err != nil {
					return nil, err
				}

				// Create a temporary Tool to get JSON Schema of sub-struct
				tempTool := mcp.NewTool("temp", subOptions...)
				tempJSON, _ := tempTool.MarshalJSON()
				var tempMap map[string]any
				if err := json.Unmarshal(tempJSON, &tempMap); err != nil {
					continue
				}

				// Extract properties from temporary Tool
				if inputSchema, ok := tempMap["inputSchema"].(map[string]any); ok {
					if properties, ok := inputSchema["properties"].(map[string]any); ok {
						// Check if there are required fields
						if required, ok := inputSchema["required"].([]any); ok {
							// Add required field information to corresponding properties
							for _, req := range required {
								if reqStr, ok := req.(string); ok {
									if prop, ok := properties[reqStr].(map[string]any); ok {
										prop["required"] = true
									}
								}
							}
						}
						items = map[string]any{
							"type":       "object",
							"properties": properties,
						}
					}
				}
			}
			opts = append(opts, mcp.Items(items))
			options = append(options, mcp.WithArray(jsonTag, opts...))
		}
	}

	return options, nil
}
