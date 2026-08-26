package api

import (
	"context"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaitin/SafeLine/mcp_server/internal/config"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

func TestMain(m *testing.M) {
	_ = logger.Init(&logger.Config{Level: "error"})
	os.Exit(m.Run())
}

func writeTokenFile(t *testing.T, token string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "api_token")
	if err := os.WriteFile(path, []byte(token+"\n"), 0o600); err != nil {
		t.Fatalf("write token file: %v", err)
	}
	return path
}

func TestClientRegistryKeepsInstanceCredentialsIsolated(t *testing.T) {
	type observedRequest struct {
		method string
		path   string
		token  string
		query  url.Values
	}
	requests := make(chan observedRequest, 2)
	newServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests <- observedRequest{
				method: r.Method,
				path:   r.URL.Path,
				token:  r.Header.Get("X-SLCE-API-TOKEN"),
				query:  r.URL.Query(),
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"nodes":[],"total":0},"err":null,"msg":""}`))
		}))
	}
	one := newServer()
	defer one.Close()
	two := newServer()
	defer two.Close()

	registry, err := NewClientRegistry([]*config.InstanceConfig{
		{ID: "dev152", DisplayName: "Development 152", BaseURL: one.URL, TokenFile: writeTokenFile(t, "token-152")},
		{ID: "dev180", DisplayName: "Development 180", BaseURL: two.URL, TokenFile: writeTokenFile(t, "token-180")},
	})
	if err != nil {
		t.Fatalf("NewClientRegistry() error = %v", err)
	}

	for _, id := range []string{"dev152", "dev180"} {
		client, err := registry.Get(id)
		if err != nil {
			t.Fatalf("Get(%q) error = %v", id, err)
		}
		var result Response[map[string]any]
		if err := client.GetAttackEvents(context.Background(), url.Values{"page": {"1"}}, &result); err != nil {
			t.Fatalf("GetAttackEvents(%q) error = %v", id, err)
		}
	}

	seenTokens := map[string]bool{}
	for range 2 {
		request := <-requests
		if request.method != http.MethodGet || request.path != attackEventsPath {
			t.Errorf("request = %s %s, want GET %s", request.method, request.path, attackEventsPath)
		}
		if request.query.Get("page") != "1" {
			t.Errorf("page query = %q", request.query.Get("page"))
		}
		seenTokens[request.token] = true
	}
	if !seenTokens["token-152"] || !seenTokens["token-180"] || len(seenTokens) != 2 {
		t.Fatalf("observed tokens = %#v, want isolated instance tokens", seenTokens)
	}

	sources := registry.Sources()
	if len(sources) != 2 || sources[0].InstanceID != "dev152" || sources[1].InstanceID != "dev180" {
		t.Fatalf("Sources() = %#v", sources)
	}
}

func TestClientReturnsHTTP200EnvelopeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{},"err":"token_invalid","msg":"invalid API token"}`))
	}))
	defer server.Close()

	client, err := newAPIClient(&config.InstanceConfig{
		ID: "broken", BaseURL: server.URL, TokenFile: writeTokenFile(t, "bad-token"),
	})
	if err != nil {
		t.Fatalf("newAPIClient() error = %v", err)
	}
	err = client.GetAttackEvents(context.Background(), nil, &Response[map[string]any]{})
	var responseErr *ResponseError
	if !stderrors.As(err, &responseErr) {
		t.Fatalf("GetAttackEvents() error = %T %v, want *ResponseError", err, err)
	}
	if responseErr.Code != "token_invalid" || responseErr.Message != "invalid API token" {
		t.Fatalf("ResponseError = %#v", responseErr)
	}
}

func TestClientRegistryRejectsUnknownInstanceAndEmptyToken(t *testing.T) {
	registry, err := NewClientRegistry([]*config.InstanceConfig{{
		ID: "known", BaseURL: "https://known.example.test", TokenFile: writeTokenFile(t, "token"),
	}})
	if err != nil {
		t.Fatalf("NewClientRegistry() error = %v", err)
	}
	if _, err := registry.Get("missing"); err == nil || !strings.Contains(err.Error(), "unknown SafeLine instance") {
		t.Fatalf("Get(missing) error = %v", err)
	}

	_, err = NewClientRegistry([]*config.InstanceConfig{{
		ID: "empty", BaseURL: "https://empty.example.test", TokenFile: writeTokenFile(t, ""),
	}})
	if err == nil || !strings.Contains(err.Error(), "token_file is empty") {
		t.Fatalf("empty token error = %v", err)
	}
}

func TestClientRejectsUnsafeTokenFiles(t *testing.T) {
	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "api_token")
	if err := os.WriteFile(tokenPath, []byte("token"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(tokenPath, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := NewClientRegistry([]*config.InstanceConfig{{
		ID: "broad-mode", BaseURL: "https://example.test", TokenFile: tokenPath,
	}})
	if err == nil || !strings.Contains(err.Error(), "no group or other permissions") {
		t.Fatalf("broad token mode error = %v", err)
	}

	if err := os.Chmod(tokenPath, 0o600); err != nil {
		t.Fatal(err)
	}
	symlinkPath := filepath.Join(dir, "token-link")
	if err := os.Symlink(tokenPath, symlinkPath); err != nil {
		t.Fatal(err)
	}
	_, err = NewClientRegistry([]*config.InstanceConfig{{
		ID: "symlink", BaseURL: "https://example.test", TokenFile: symlinkPath,
	}})
	if err == nil || !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("symlink token error = %v", err)
	}
}

func TestClientLimitsSafeLineResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", maxResponseBodySize+1)))
	}))
	defer server.Close()

	client, err := newAPIClient(&config.InstanceConfig{
		ID: "oversized", BaseURL: server.URL, TokenFile: writeTokenFile(t, "token"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.GetAttackEvents(context.Background(), nil, &Response[map[string]any]{})
	if err == nil || !strings.Contains(err.Error(), "4 MiB limit") {
		t.Fatalf("oversized response error = %v", err)
	}
}
