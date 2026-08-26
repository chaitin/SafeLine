package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaitin/SafeLine/mcp_server/internal/auth"
	"github.com/chaitin/SafeLine/mcp_server/internal/config"
)

func TestAuthInitIsOneTimeAndRotateIsExplicit(t *testing.T) {
	stateFile := filepath.Join(t.TempDir(), "state", "auth.json")
	var first bytes.Buffer
	if err := run([]string{"auth", "init", "--state-file", stateFile}, &first, &bytes.Buffer{}); err != nil {
		t.Fatalf("auth init error = %v", err)
	}
	firstToken := lastLine(first.String())
	if !strings.HasPrefix(firstToken, "slmcp_") {
		t.Fatalf("init output did not contain generated token: %q", first.String())
	}

	contents, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(contents), firstToken) {
		t.Fatal("state file contains plaintext token")
	}

	var second bytes.Buffer
	if err := run([]string{"auth", "init", "--state-file", stateFile}, &second, &bytes.Buffer{}); err != nil {
		t.Fatalf("second auth init error = %v", err)
	}
	if strings.Contains(second.String(), "slmcp_") {
		t.Fatalf("second init exposed or generated a token: %q", second.String())
	}

	var rotated bytes.Buffer
	if err := run([]string{"auth", "rotate", "--state-file", stateFile}, &rotated, &bytes.Buffer{}); err != nil {
		t.Fatalf("auth rotate error = %v", err)
	}
	rotatedToken := lastLine(rotated.String())
	if rotatedToken == firstToken || !strings.HasPrefix(rotatedToken, "slmcp_") {
		t.Fatalf("rotate output token = %q", rotatedToken)
	}
	a, err := auth.Load(stateFile)
	if err != nil {
		t.Fatal(err)
	}
	if a.Verify(firstToken) || !a.Verify(rotatedToken) {
		t.Fatal("explicit rotation did not replace the accepted token")
	}
}

func lastLine(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	return lines[len(lines)-1]
}

func TestBuildServerInstructionsIncludesFriendlyInstanceMappings(t *testing.T) {
	instructions := buildServerInstructions([]*config.InstanceConfig{
		{
			ID:          "dev152",
			DisplayName: "Development 152",
			BaseURL:     "https://secret.internal",
			TokenFile:   "/run/secrets/dev152.token",
		},
		{
			ID:          "dev180",
			DisplayName: "Development 180",
			BaseURL:     "https://other.internal",
			TokenFile:   "/run/secrets/dev180.token",
		},
	})

	for _, expected := range []string{
		`display_name "Development 152" -> instance_id "dev152"`,
		`display_name "Development 180" -> instance_id "dev180"`,
		"Never guess an instance_id",
	} {
		if !strings.Contains(instructions, expected) {
			t.Errorf("instructions = %q, want substring %q", instructions, expected)
		}
	}
	for _, sensitive := range []string{
		"secret.internal",
		"other.internal",
		"/run/secrets/dev152.token",
		"/run/secrets/dev180.token",
	} {
		if strings.Contains(instructions, sensitive) {
			t.Fatalf("instructions exposed instance configuration %q: %q", sensitive, instructions)
		}
	}
}
