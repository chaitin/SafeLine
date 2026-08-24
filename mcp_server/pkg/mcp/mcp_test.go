package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	markmcp "github.com/mark3labs/mcp-go/mcp"
)

type probeParams struct {
	Count int `json:"count" jsonschema:"Number of results to return"`
}

type probeResult struct {
	Count int `json:"count" jsonschema:"Returned result count"`
}

type probeTool struct{}

func (probeTool) Name() string        { return "probe" }
func (probeTool) Description() string { return "Protocol test tool" }
func (probeTool) Validate(params probeParams) error {
	if params.Count < 1 {
		return errors.New("count must be at least 1")
	}
	return nil
}
func (probeTool) Execute(_ context.Context, params probeParams) (probeResult, error) {
	if params.Count == 13 {
		return probeResult{}, errors.New("downstream failed")
	}
	return probeResult{Count: params.Count}, nil
}

func TestModernStreamableHTTPAndToolErrorSemantics(t *testing.T) {
	if err := logger.Init(&logger.Config{Level: "error"}); err != nil {
		t.Fatal(err)
	}
	s := NewMCPServer("test", "1.0.0")
	if err := RegisterTool[probeParams, probeResult](s, probeTool{}); err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(s.Handler())
	defer httpServer.Close()

	discover := postModern(t, httpServer.URL+EndpointPath, markmcp.MethodServerDiscover, map[string]any{})
	result := requireRPCResult(t, discover)
	versions, ok := result["supportedVersions"].([]any)
	if !ok || len(versions) == 0 || versions[0] != markmcp.ProtocolVersion20260728 {
		t.Fatalf("supportedVersions = %#v, want %s first", result["supportedVersions"], markmcp.ProtocolVersion20260728)
	}

	valid := postModern(t, httpServer.URL+EndpointPath, markmcp.MethodToolsCall, map[string]any{
		"name":      "probe",
		"arguments": map[string]any{"count": 2},
	})
	validResult := requireRPCResult(t, valid)
	if validResult["isError"] == true {
		t.Fatalf("valid tool result = %#v", validResult)
	}
	structured := validResult["structuredContent"].(map[string]any)
	if structured["count"] != float64(2) {
		t.Fatalf("structuredContent = %#v", structured)
	}

	for name, arguments := range map[string]map[string]any{
		"decimal integer": {"count": 1.5},
		"unknown field":   {"count": 1, "unexpected": true},
		"validation":      {"count": 0},
		"downstream":      {"count": 13},
	} {
		t.Run(name, func(t *testing.T) {
			response := postModern(t, httpServer.URL+EndpointPath, markmcp.MethodToolsCall, map[string]any{
				"name":      "probe",
				"arguments": arguments,
			})
			toolResult := requireRPCResult(t, response)
			if toolResult["isError"] != true {
				t.Fatalf("tool result = %#v, want isError=true", toolResult)
			}
		})
	}

	for _, path := range []string{"/sse", "/message"} {
		response, err := http.Get(httpServer.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("GET %s status = %d, want 404", path, response.StatusCode)
		}
	}
}

func TestHandlerRejectsOversizedRequestBody(t *testing.T) {
	s := NewMCPServer("test", "1.0.0")
	body := strings.NewReader(strings.Repeat("x", maxRequestBodySize+1))
	req := httptest.NewRequest(http.MethodPost, EndpointPath, body)
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	s.Handler().ServeHTTP(response, req)
	if !strings.Contains(response.Body.String(), "request body too large") {
		t.Fatalf("oversized response = status %d body %q", response.Code, response.Body.String())
	}
}

func postModern(t *testing.T, url string, method markmcp.MCPMethod, params map[string]any) *http.Response {
	t.Helper()
	params["_meta"] = map[string]any{
		markmcp.MetaKeyProtocolVersion: markmcp.ProtocolVersion20260728,
		markmcp.MetaKeyClientInfo: map[string]any{
			"name":    "safeline-mcp-test",
			"version": "1.0.0",
		},
		markmcp.MetaKeyClientCapabilities: map[string]any{},
	}
	body, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set(markmcp.HeaderProtocolVersion, markmcp.ProtocolVersion20260728)
	req.Header.Set(markmcp.HeaderMethod, string(method))
	if name, ok := params["name"].(string); ok {
		encoded, valid := markmcp.EncodeHeaderValue(name)
		if !valid {
			t.Fatalf("invalid MCP name %q", name)
		}
		req.Header.Set(markmcp.HeaderName, encoded)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = response.Body.Close() })
	if response.StatusCode != http.StatusOK {
		var payload any
		_ = json.NewDecoder(response.Body).Decode(&payload)
		t.Fatalf("POST %s status = %d, payload = %#v", method, response.StatusCode, payload)
	}
	return response
}

func requireRPCResult(t *testing.T, response *http.Response) map[string]any {
	t.Helper()
	var rpc map[string]any
	if err := json.NewDecoder(response.Body).Decode(&rpc); err != nil {
		t.Fatal(err)
	}
	result, ok := rpc["result"].(map[string]any)
	if !ok {
		t.Fatalf("JSON-RPC response = %#v", rpc)
	}
	if result["resultType"] != string(markmcp.ResultTypeComplete) {
		t.Fatalf("resultType = %#v, want %q", result["resultType"], markmcp.ResultTypeComplete)
	}
	return result
}
