package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfigFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadMultiInstanceConfiguration(t *testing.T) {
	path := writeConfigFile(t, `
server:
  name: SafeLine MCP
  version: 1.0.0
  host: 127.0.0.1
  port: 5678
logger:
  level: info
instances:
  - id: dev152
    display_name: Development 152
    base_url: https://dev152:9443
    token_file: secrets/dev152_token
    timeout: 12
    insecure_skip_verify: true
  - id: dev180
    base_url: https://dev180:9443
    token_file: /run/secrets/dev180_token
`)

	if err := Load(path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	instances := GetInstances()
	if len(instances) != 2 {
		t.Fatalf("len(GetInstances()) = %d, want 2", len(instances))
	}
	if instances[0].ID != "dev152" || instances[0].DisplayName != "Development 152" {
		t.Fatalf("first instance = %#v", instances[0])
	}
	wantRelativePath := filepath.Join(filepath.Dir(path), "secrets", "dev152_token")
	if instances[0].TokenFile != wantRelativePath {
		t.Errorf("relative token_file = %q, want %q", instances[0].TokenFile, wantRelativePath)
	}
	if instances[1].DisplayName != "dev180" {
		t.Errorf("default display_name = %q, want dev180", instances[1].DisplayName)
	}
	if instances[1].TokenFile != "/run/secrets/dev180_token" {
		t.Errorf("absolute token_file = %q", instances[1].TokenFile)
	}
}

func TestLoadRejectsDuplicateInstanceIDs(t *testing.T) {
	path := writeConfigFile(t, `
server: {name: SafeLine MCP, version: 1.0.0, host: 127.0.0.1, port: 5678}
logger: {level: info}
instances:
  - id: duplicate
    base_url: https://one.example.test
    token_file: one_token
  - id: duplicate
    base_url: https://two.example.test
    token_file: two_token
`)

	err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "duplicate instance id") {
		t.Fatalf("Load() error = %v, want duplicate instance id", err)
	}
}

func TestLoadRejectsDuplicateDisplayNames(t *testing.T) {
	path := writeConfigFile(t, `
server: {name: SafeLine MCP, version: 1.0.0, host: 127.0.0.1, port: 5678}
logger: {level: info}
instances:
  - id: production-a
    display_name: Production
    base_url: https://one.example.test
    token_file: one_token
  - id: production-b
    display_name: " production "
    base_url: https://two.example.test
    token_file: two_token
`)

	err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "duplicate instance display_name") {
		t.Fatalf("Load() error = %v, want duplicate instance display_name", err)
	}
}

func TestLoadRejectsDisplayNameConflictingWithAnotherInstanceID(t *testing.T) {
	path := writeConfigFile(t, `
server: {name: SafeLine MCP, version: 1.0.0, host: 127.0.0.1, port: 5678}
logger: {level: info}
instances:
  - id: production-a
    display_name: PRODUCTION-B
    base_url: https://one.example.test
    token_file: one_token
  - id: production-b
    display_name: Backup Production
    base_url: https://two.example.test
    token_file: two_token
`)

	err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "conflicts with instance id") {
		t.Fatalf("Load() error = %v, want display_name and instance id conflict", err)
	}
}

func TestLoadRejectsInlineInstanceToken(t *testing.T) {
	path := writeConfigFile(t, `
server: {name: SafeLine MCP, version: 1.0.0, host: 127.0.0.1, port: 5678}
logger: {level: info}
instances:
  - id: unsafe
    base_url: https://unsafe.example.test
    token: must-not-be-accepted
`)

	err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "field token not found") {
		t.Fatalf("Load() error = %v, want unknown inline token field", err)
	}
}

func TestLoadRequiresAtLeastOneInstance(t *testing.T) {
	path := writeConfigFile(t, `
server: {name: SafeLine MCP, version: 1.0.0, host: 127.0.0.1, port: 5678}
logger: {level: info}
`)
	err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "at least one SafeLine instance") {
		t.Fatalf("Load() error = %v, want missing instance error", err)
	}
}
