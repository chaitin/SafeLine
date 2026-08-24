package analyze

import (
	"encoding/json"
	"testing"

	markmcp "github.com/mark3labs/mcp-go/mcp"
)

func TestGetAttackEventsValidate(t *testing.T) {
	tool := &GetAttackEvents{}
	valid := GetAttackEventsParams{InstanceID: "dev152", Page: 1, PageSize: 10}
	if err := tool.Validate(valid); err != nil {
		t.Fatalf("Validate(valid) error = %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*GetAttackEventsParams)
	}{
		{name: "missing instance", mutate: func(p *GetAttackEventsParams) { p.InstanceID = " " }},
		{name: "page zero", mutate: func(p *GetAttackEventsParams) { p.Page = 0 }},
		{name: "page size zero", mutate: func(p *GetAttackEventsParams) { p.PageSize = 0 }},
		{name: "page size too large", mutate: func(p *GetAttackEventsParams) { p.PageSize = 101 }},
		{name: "negative port", mutate: func(p *GetAttackEventsParams) { p.Port = -1 }},
		{name: "port too large", mutate: func(p *GetAttackEventsParams) { p.Port = 65536 }},
		{name: "negative start", mutate: func(p *GetAttackEventsParams) { p.Start = -1 }},
		{name: "negative end", mutate: func(p *GetAttackEventsParams) { p.End = -1 }},
		{name: "reversed range", mutate: func(p *GetAttackEventsParams) { p.Start, p.End = 2, 1 }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			params := valid
			test.mutate(&params)
			if err := tool.Validate(params); err == nil {
				t.Fatal("Validate() error = nil, want validation error")
			}
		})
	}
}

func TestGetAttackEventsV1InputSchema(t *testing.T) {
	implementation := &GetAttackEvents{}
	tool := markmcp.NewTool("get_attack_events", implementation.InputSchemaOptions()...)
	var schema struct {
		Properties           map[string]map[string]any `json:"properties"`
		Required             []string                  `json:"required"`
		AdditionalProperties bool                      `json:"additionalProperties"`
	}
	encoded, err := json.Marshal(tool.InputSchema)
	if err != nil {
		t.Fatalf("marshal input schema: %v", err)
	}
	if err := json.Unmarshal(encoded, &schema); err != nil {
		t.Fatalf("unmarshal input schema: %v", err)
	}
	if len(schema.Required) != 1 || schema.Required[0] != "instance_id" {
		t.Fatalf("required = %#v, want only instance_id", schema.Required)
	}
	for _, field := range []string{"page", "page_size", "port", "start", "end"} {
		if got := schema.Properties[field]["type"]; got != "integer" {
			t.Errorf("%s type = %v, want integer", field, got)
		}
	}
	if got := schema.Properties["page"]["minimum"]; got != float64(1) {
		t.Errorf("page minimum = %v, want 1", got)
	}
	if got := schema.Properties["page_size"]["maximum"]; got != float64(100) {
		t.Errorf("page_size maximum = %v, want 100", got)
	}
	if got := schema.Properties["port"]["maximum"]; got != float64(65535) {
		t.Errorf("port maximum = %v, want 65535", got)
	}
	if schema.AdditionalProperties {
		t.Error("additionalProperties = true, want false")
	}
}
