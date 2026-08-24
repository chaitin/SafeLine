package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaitin/SafeLine/mcp_server/internal/auth"
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
