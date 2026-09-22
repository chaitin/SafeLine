package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureBootstrapTokenCreatesAndReusesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tfa_bootstrap_token")
	t.Setenv(BootstrapTokenEnv, "")
	t.Setenv(BootstrapTokenFileEnv, path)

	created, gotPath, err := EnsureBootstrapToken()
	if err != nil {
		t.Fatal(err)
	}
	if !created || gotPath != path || bootstrapToken == "" {
		t.Fatalf("created=%v path=%q token empty=%v", created, gotPath, bootstrapToken == "")
	}
	first := bootstrapToken

	created, _, err = EnsureBootstrapToken()
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("second call rotated the token")
	}
	if bootstrapToken != first {
		t.Fatal("token changed")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) == "" {
		t.Fatal("file empty")
	}
}

func TestEnsureBootstrapTokenPrefersEnv(t *testing.T) {
	t.Setenv(BootstrapTokenEnv, "from-env")
	t.Setenv(BootstrapTokenFileEnv, filepath.Join(t.TempDir(), "unused"))
	created, path, err := EnsureBootstrapToken()
	if err != nil {
		t.Fatal(err)
	}
	if created || path != "" || bootstrapToken != "from-env" {
		t.Fatalf("created=%v path=%q token=%q", created, path, bootstrapToken)
	}
}

func TestEnsureBootstrapTokenReplacesEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tfa_bootstrap_token")
	if err := os.WriteFile(path, []byte(" \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(BootstrapTokenEnv, "")
	t.Setenv(BootstrapTokenFileEnv, path)
	created, _, err := EnsureBootstrapToken()
	if err != nil {
		t.Fatal(err)
	}
	if !created || bootstrapToken == "" {
		t.Fatalf("created=%v token empty", created)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(content)) != bootstrapToken {
		t.Fatalf("file=%q token=%q", content, bootstrapToken)
	}
}
