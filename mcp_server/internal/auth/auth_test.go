package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestInitializePersistsOnlyHashAndKeepsTokenAcrossLoads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "auth.json")

	token, created, err := Initialize(path)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if !created || !strings.HasPrefix(token, "slmcp_") {
		t.Fatalf("Initialize() = (%q, %v), want newly generated prefixed token", token, created)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(contents), token) {
		t.Fatal("authentication state contains the plaintext token")
	}
	var persisted state
	if err := json.Unmarshal(contents, &persisted); err != nil {
		t.Fatal(err)
	}
	if persisted.TokenHash == "" {
		t.Fatal("authentication state does not contain a token hash")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("state mode = %o, want 600", got)
	}

	a, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !a.Verify(token) || a.Verify(token+"wrong") {
		t.Fatal("loaded authenticator did not verify exactly the initialized token")
	}

	second, secondCreated, err := Initialize(path)
	if err != nil {
		t.Fatalf("second Initialize() error = %v", err)
	}
	if second != "" || secondCreated {
		t.Fatalf("second Initialize() = (%q, %v), want no rotation", second, secondCreated)
	}
}

func TestRotateInvalidatesPreviousToken(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "auth.json")
	oldToken, _, err := Initialize(path)
	if err != nil {
		t.Fatal(err)
	}
	newToken, err := Rotate(path)
	if err != nil {
		t.Fatal(err)
	}
	if oldToken == newToken {
		t.Fatal("Rotate() returned the previous token")
	}

	a, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if a.Verify(oldToken) || !a.Verify(newToken) {
		t.Fatal("rotation did not replace the accepted token")
	}
}

func TestInitializeIsAtomicAcrossConcurrentCallers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "auth.json")
	const callers = 64

	type initializeResult struct {
		token   string
		created bool
		err     error
	}
	start := make(chan struct{})
	results := make(chan initializeResult, callers)
	var group sync.WaitGroup
	for range callers {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			token, created, err := Initialize(path)
			results <- initializeResult{token: token, created: created, err: err}
		}()
	}
	close(start)
	group.Wait()
	close(results)

	var winningToken string
	createdCount := 0
	for result := range results {
		if result.err != nil {
			t.Fatalf("Initialize() error = %v", result.err)
		}
		if result.created {
			createdCount++
			winningToken = result.token
		} else if result.token != "" {
			t.Fatalf("losing initializer exposed token %q", result.token)
		}
	}
	if createdCount != 1 {
		t.Fatalf("created count = %d, want exactly 1", createdCount)
	}
	a, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Verify(winningToken) {
		t.Fatal("persisted state does not accept the one exposed token")
	}
}

func TestInitializeRejectsExistingInvalidState(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "auth.json")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if token, created, err := Initialize(path); err == nil {
		t.Fatalf("Initialize() = (%q, %v, nil), want invalid-state error", token, created)
	}
}

func TestRotateRequiresExistingValidState(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "missing.json")
	if token, err := Rotate(path); err == nil {
		t.Fatalf("Rotate() = (%q, nil), want missing-state error", token)
	}

	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if token, err := Rotate(path); err == nil {
		t.Fatalf("Rotate(invalid) = (%q, nil), want invalid-state error", token)
	}
}

func TestLoadRejectsOverlyBroadStatePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "auth.json")
	if _, _, err := Initialize(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load() accepted authentication state readable by other users")
	}
}

func TestLoadRejectsOversizedStateFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "auth.json")
	if err := os.WriteFile(path, []byte(strings.Repeat("x", maxStateFileSize+1)), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "16 KiB limit") {
		t.Fatalf("Load(oversized) error = %v", err)
	}
}

func TestInitializeRejectsOverlyBroadStateDirectoryPermissions(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "auth.json")
	if _, _, err := Initialize(path); err == nil {
		t.Fatal("Initialize() accepted a state directory accessible by group or other users")
	}
}

func TestMiddlewareRequiresBearerToken(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "auth.json")
	token, _, err := Initialize(path)
	if err != nil {
		t.Fatal(err)
	}
	a, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	handler := a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for name, authorization := range map[string]string{
		"missing":      "",
		"wrong scheme": "Basic " + token,
		"wrong token":  "Bearer wrong",
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
			req.Header.Set("Authorization", authorization)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", response.Code)
			}
		})
	}

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", response.Code)
	}
}

func TestEnabledFromEnv(t *testing.T) {
	t.Setenv(EnabledEnv, "true")
	enabled, err := EnabledFromEnv()
	if err != nil || !enabled {
		t.Fatalf("EnabledFromEnv() = (%v, %v), want true", enabled, err)
	}

	t.Setenv(EnabledEnv, "not-a-bool")
	if _, err := EnabledFromEnv(); err == nil {
		t.Fatal("EnabledFromEnv() accepted an invalid boolean")
	}
}
