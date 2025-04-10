package mcp

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestSchemaToOptions(t *testing.T) {
	tests := []struct {
		name string
		args any
		want mcp.Tool
	}{
		{
			name: "test number",
			args: struct {
				A int `json:"a" desc:"number a" required:"true"`
			}{},
			want: mcp.NewTool("test number",
				mcp.WithNumber("a", mcp.Required(), mcp.Description("number a")),
			),
		},
		{
			name: "test number int64",
			args: struct {
				A int64 `json:"a" desc:"number a" required:"true"`
			}{},
			want: mcp.NewTool("test number int64",
				mcp.WithNumber("a", mcp.Required(), mcp.Description("number a")),
			),
		},
		{
			name: "test number float64",
			args: struct {
				A float64 `json:"a" desc:"number a" required:"true"`
			}{},
			want: mcp.NewTool("test number float64",
				mcp.WithNumber("a", mcp.Required(), mcp.Description("number a")),
			),
		},
		{
			name: "test number default",
			args: struct {
				A int `json:"a" desc:"number a" required:"true" default:"10"`
			}{},
			want: mcp.NewTool("test number default",
				mcp.WithNumber("a", mcp.Required(), mcp.Description("number a"), mcp.DefaultNumber(10)),
			),
		},
		{
			name: "test number min max",
			args: struct {
				A int `json:"a" desc:"number a" required:"true" min:"10" max:"20"`
			}{},
			want: mcp.NewTool("test number min max",
				mcp.WithNumber("a", mcp.Required(), mcp.Description("number a"), mcp.Min(10), mcp.Max(20)),
			),
		},
		{
			name: "test number optional",
			args: struct {
				A int `json:"a" desc:"number a"`
			}{},
			want: mcp.NewTool("test number optional",
				mcp.WithNumber("a", mcp.Description("number a")),
			),
		},
		{
			name: "test boolean",
			args: struct {
				A bool `json:"a" desc:"boolean a" required:"true"`
			}{},
			want: mcp.NewTool("test boolean",
				mcp.WithBoolean("a", mcp.Required(), mcp.Description("boolean a")),
			),
		},
		{
			name: "test string",
			args: struct {
				A string `json:"a" desc:"string a" required:"true"`
			}{},
			want: mcp.NewTool("test string",
				mcp.WithString("a", mcp.Required(), mcp.Description("string a")),
			),
		},
		{
			name: "test string default",
			args: struct {
				A string `json:"a" desc:"string a" required:"true" default:"hello"`
			}{},
			want: mcp.NewTool("test string default",
				mcp.WithString("a", mcp.Required(), mcp.Description("string a"), mcp.DefaultString("hello")),
			),
		},
		{
			name: "test string enum",
			args: struct {
				A string `json:"a" desc:"string a" required:"true" enum:"1,2,3"`
			}{},
			want: mcp.NewTool("test string enum",
				mcp.WithString("a", mcp.Required(), mcp.Description("string a"), mcp.Enum("1", "2", "3")),
			),
		},
		{
			name: "test object",
			args: struct {
				A struct {
					B int `json:"b" desc:"number b" required:"true"`
				} `json:"a" desc:"object a" required:"true"`
			}{},
			want: mcp.NewTool("test object",
				mcp.WithObject("a", mcp.Required(), mcp.Description("object a"),
					mcp.Properties(map[string]any{
						"b": map[string]any{
							"type":        "number",
							"description": "number b",
							"required":    true,
						},
					}),
				),
			),
		},
		{
			name: "test object optional",
			args: struct {
				A struct {
					B int `json:"b" desc:"number b" required:"true"`
					C int `json:"c" desc:"number c"`
				} `json:"a" desc:"object a" required:"true"`
			}{},
			want: mcp.NewTool("test object optional",
				mcp.WithObject("a", mcp.Required(), mcp.Description("object a"),
					mcp.Properties(map[string]any{
						"b": map[string]any{
							"type":        "number",
							"description": "number b",
							"required":    true,
						},
						"c": map[string]any{
							"type":        "number",
							"description": "number c",
						},
					}),
				),
			),
		},
		{
			name: "test nested object",
			args: struct {
				A struct {
					B struct {
						C int `json:"c" desc:"number c" required:"true"`
					} `json:"b" desc:"object b" required:"true"`
				} `json:"a" desc:"object a" required:"true"`
			}{},
			want: mcp.NewTool("test nested object",
				mcp.WithObject("a", mcp.Required(), mcp.Description("object a"),
					mcp.Properties(map[string]any{
						"b": map[string]any{
							"type":        "object",
							"description": "object b",
							"required":    true,
							"properties": map[string]any{
								"c": map[string]any{
									"type":        "number",
									"description": "number c",
									"required":    true,
								},
							},
						},
					}),
				),
			),
		},
		{
			name: "test array",
			args: struct {
				A []int `json:"a" desc:"array a" required:"true"`
			}{},
			want: mcp.NewTool("test array",
				mcp.WithArray("a", mcp.Required(), mcp.Description("array a"),
					mcp.Items(map[string]any{
						"type": "number",
					}),
				),
			),
		},
		{
			name: "test array of object",
			args: struct {
				A []struct {
					B int `json:"b" desc:"number b" required:"true"`
				} `json:"a" desc:"array of object a" required:"true"`
			}{},
			want: mcp.NewTool("test array of object",
				mcp.WithArray("a", mcp.Required(), mcp.Description("array of object a"),
					mcp.Items(map[string]any{
						"type": "object",
						"properties": map[string]any{
							"b": map[string]any{
								"type":        "number",
								"description": "number b",
								"required":    true,
							},
						},
					}),
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SchemaToOptions(tt.args)
			if err != nil {
				t.Errorf("SchemaToOptions() error = %v", err)
				return
			}
			s1, _ := mcp.NewTool(tt.name, got...).MarshalJSON()
			s2, _ := tt.want.MarshalJSON()
			if !reflect.DeepEqual(s1, s2) {
				t.Errorf("\n got  %v\n want %v", string(s1), string(s2))
			}
		})
	}
}
