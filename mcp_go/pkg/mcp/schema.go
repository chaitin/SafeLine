package mcp

import (
	"encoding/json"
	"reflect"
	"strconv"
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
		defaultTag := field.Tag.Get("default")
		minTag := field.Tag.Get("min")
		maxTag := field.Tag.Get("max")
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
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			if defaultTag != "" {
				if defaultValue, err := strconv.Atoi(defaultTag); err == nil {
					opts = append(opts, mcp.DefaultNumber(float64(defaultValue)))
				}
			}
			if minTag != "" {
				if minValue, err := strconv.Atoi(minTag); err == nil {
					opts = append(opts, mcp.Min(float64(minValue)))
				}
			}
			if maxTag != "" {
				if maxValue, err := strconv.Atoi(maxTag); err == nil {
					opts = append(opts, mcp.Max(float64(maxValue)))
				}
			}
			options = append(options, mcp.WithNumber(jsonTag, opts...))
		case reflect.Bool:
			if defaultTag != "" {
				if defaultValue, err := strconv.ParseBool(defaultTag); err == nil {
					opts = append(opts, mcp.DefaultBool(defaultValue))
				}
			}
			options = append(options, mcp.WithBoolean(jsonTag, opts...))
		case reflect.String:
			if defaultTag != "" {
				opts = append(opts, mcp.DefaultString(defaultTag))
			}
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
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
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
