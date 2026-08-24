package analyze

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/internal/config"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

func TestGetEventListUsesCurrentSafeLineQueryAndReturnsSource(t *testing.T) {
	if err := logger.Init(&logger.Config{Level: "error"}); err != nil {
		t.Fatal(err)
	}
	tokenPath := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tokenPath, []byte("instance-token\n"), 0o600); err != nil {
		t.Fatalf("write token: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/open/events" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("X-SLCE-API-TOKEN"); got != "instance-token" {
			t.Errorf("token header = %q", got)
		}
		wantQuery := map[string]string{
			"page": "2", "page_size": "25", "ip": "10.0.0.0/24",
			"host": "example.test&a=b", "port": "9443", "start": "1000", "end": "2000",
		}
		for key, want := range wantQuery {
			if got := r.URL.Query().Get(key); got != want {
				t.Errorf("query[%q] = %q, want %q", key, got, want)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
          "data": {
            "nodes": [{
              "id": 7,
              "ip": "10.0.0.8",
              "protocol": 1,
              "host": "example.test",
              "dst_port": 9443,
              "updated_at": 1800,
              "start_at": 1100,
              "end_at": 1700,
              "deny_count": 3,
              "pass_count": 0,
              "finished": true,
              "country": "CN",
              "province": "Beijing",
              "city": "Beijing"
            }],
            "total": 1
          },
          "err": null,
          "msg": ""
        }`))
	}))
	defer server.Close()

	if err := api.InitInstances([]*config.InstanceConfig{{
		ID: "dev152", DisplayName: "Development 152", BaseURL: server.URL, TokenFile: tokenPath,
	}}); err != nil {
		t.Fatalf("InitInstances() error = %v", err)
	}

	response, err := GetEventList(context.Background(), &GetEventListRequest{
		InstanceID: "dev152",
		Page:       2,
		PageSize:   25,
		IP:         "10.0.0.0/24",
		Host:       "example.test&a=b",
		Port:       9443,
		Start:      1000,
		End:        2000,
	})
	if err != nil {
		t.Fatalf("GetEventList() error = %v", err)
	}
	if response.Source.InstanceID != "dev152" || response.Source.DisplayName != "Development 152" {
		t.Errorf("source = %#v", response.Source)
	}
	if response.Page != 2 || response.PageSize != 25 || response.Total != 1 {
		t.Errorf("pagination/total = page %d, size %d, total %d", response.Page, response.PageSize, response.Total)
	}
	if len(response.Nodes) != 1 || response.Nodes[0].ID != 7 || response.Nodes[0].DstPort != 9443 {
		t.Fatalf("nodes = %#v", response.Nodes)
	}
}
