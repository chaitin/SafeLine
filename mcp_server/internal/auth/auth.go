package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	EnabledEnv   = "MCP_AUTH_ENABLED"
	StateFileEnv = "MCP_AUTH_STATE_FILE"

	DefaultStateFile = "/var/lib/safeline-mcp/auth.json"
	stateVersion     = 1
	tokenBytes       = 32
	maxStateFileSize = 16 << 10
)

var ErrInvalidState = errors.New("invalid MCP authentication state")

type state struct {
	Version   int       `json:"version"`
	TokenHash string    `json:"token_hash"`
	CreatedAt time.Time `json:"created_at"`
}

// Authenticator verifies the deployment-scoped bearer token without retaining
// the plaintext token after initialization.
type Authenticator struct {
	tokenHash [sha256.Size]byte
}

func EnabledFromEnv() (bool, error) {
	raw := strings.TrimSpace(os.Getenv(EnabledEnv))
	if raw == "" {
		return false, nil
	}

	enabled, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("parse %s: %w", EnabledEnv, err)
	}
	return enabled, nil
}

func StateFileFromEnv() string {
	if path := strings.TrimSpace(os.Getenv(StateFileEnv)); path != "" {
		return path
	}
	return DefaultStateFile
}

// Initialize creates the authentication state once. The returned plaintext
// token must be shown to the deployer exactly once and is never written to disk.
func Initialize(path string) (token string, created bool, err error) {
	if strings.TrimSpace(path) == "" {
		return "", false, errors.New("authentication state file path is required")
	}

	if _, err := os.Lstat(path); err == nil {
		if _, err := Load(path); err != nil {
			return "", false, fmt.Errorf("existing authentication state is invalid: %w", err)
		}
		return "", false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", false, fmt.Errorf("check authentication state: %w", err)
	}

	token, s, err := newState()
	if err != nil {
		return "", false, err
	}
	created, err = writeState(path, s, false)
	if err != nil {
		return "", false, err
	}
	if !created {
		// Another initializer won the no-replace install race. Validate its
		// fully synced state and never expose this process's discarded token.
		if _, err := Load(path); err != nil {
			return "", false, fmt.Errorf("concurrently created authentication state is invalid: %w", err)
		}
		return "", false, nil
	}
	return token, true, nil
}

// Rotate replaces an existing, valid token state. Initialization is an
// explicit separate operation so a typo cannot overwrite an unrelated path.
func Rotate(path string) (string, error) {
	if _, err := Load(path); err != nil {
		return "", fmt.Errorf("load existing authentication state before rotation: %w", err)
	}
	token, s, err := newState()
	if err != nil {
		return "", err
	}
	if _, err := writeState(path, s, true); err != nil {
		return "", err
	}
	return token, nil
}

func newState() (token string, s state, err error) {
	token, err = generateToken()
	if err != nil {
		return "", state{}, err
	}

	digest := sha256.Sum256([]byte(token))
	now := time.Now().UTC()
	s = state{
		Version:   stateVersion,
		TokenHash: hex.EncodeToString(digest[:]),
		CreatedAt: now,
	}
	return token, s, nil
}

func generateToken() (string, error) {
	random := make([]byte, tokenBytes)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate bearer token: %w", err)
	}
	return "slmcp_" + base64.RawURLEncoding.EncodeToString(random), nil
}

// writeState writes a fully synced temporary file before installation. When
// replace is false, a hard-link install provides atomic create-if-absent
// semantics: concurrent initializers cannot overwrite the winning state.
func writeState(path string, s state, replace bool) (bool, error) {
	if strings.TrimSpace(path) == "" {
		return false, errors.New("authentication state file path is required")
	}

	dir := filepath.Dir(path)
	info, err := os.Lstat(dir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return false, fmt.Errorf("create authentication state directory: %w", err)
		}
		if err := os.Chmod(dir, 0o700); err != nil {
			return false, fmt.Errorf("secure authentication state directory: %w", err)
		}
		info, err = os.Lstat(dir)
	} else if err != nil {
		return false, fmt.Errorf("inspect authentication state directory: %w", err)
	}
	if err != nil {
		return false, fmt.Errorf("inspect authentication state directory: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return false, errors.New("authentication state directory path must be a real directory")
	}
	if info.Mode().Perm()&0o077 != 0 {
		return false, errors.New("authentication state directory must have no group or other permissions")
	}

	tmp, err := os.CreateTemp(dir, ".auth-*.tmp")
	if err != nil {
		return false, fmt.Errorf("create authentication state: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	if err := tmp.Chmod(0o600); err != nil {
		return false, fmt.Errorf("secure authentication state: %w", err)
	}
	encoder := json.NewEncoder(tmp)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(s); err != nil {
		return false, fmt.Errorf("encode authentication state: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		return false, fmt.Errorf("sync authentication state: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return false, fmt.Errorf("close authentication state: %w", err)
	}

	if replace {
		if err := os.Rename(tmpPath, path); err != nil {
			return false, fmt.Errorf("install authentication state: %w", err)
		}
		return true, nil
	}

	if err := os.Link(tmpPath, path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return false, nil
		}
		return false, fmt.Errorf("install authentication state without replacement: %w", err)
	}
	return true, nil
}

func Load(path string) (*Authenticator, error) {
	dirInfo, err := os.Lstat(filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("inspect authentication state directory: %w", err)
	}
	if dirInfo.Mode()&os.ModeSymlink != 0 || !dirInfo.IsDir() || dirInfo.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("%w: state directory must be a real directory with no group or other permissions", ErrInvalidState)
	}

	pathInfo, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect authentication state: %w", err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 || !pathInfo.Mode().IsRegular() || pathInfo.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("%w: state file must be a regular file with no group or other permissions", ErrInvalidState)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open authentication state: %w", err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("inspect opened authentication state: %w", err)
	}
	if !os.SameFile(pathInfo, openedInfo) || !openedInfo.Mode().IsRegular() || openedInfo.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("%w: authentication state changed while opening", ErrInvalidState)
	}
	contents, err := io.ReadAll(io.LimitReader(file, maxStateFileSize+1))
	if err != nil {
		return nil, fmt.Errorf("read authentication state: %w", err)
	}
	if len(contents) > maxStateFileSize {
		return nil, fmt.Errorf("%w: state file exceeds the 16 KiB limit", ErrInvalidState)
	}

	var s state
	if err := json.Unmarshal(contents, &s); err != nil {
		return nil, fmt.Errorf("%w: decode state: %v", ErrInvalidState, err)
	}
	if s.Version != stateVersion || s.CreatedAt.IsZero() {
		return nil, ErrInvalidState
	}
	digest, err := hex.DecodeString(s.TokenHash)
	if err != nil || len(digest) != sha256.Size {
		return nil, ErrInvalidState
	}

	a := &Authenticator{}
	copy(a.tokenHash[:], digest)
	return a, nil
}

func (a *Authenticator) Verify(token string) bool {
	if a == nil || token == "" {
		return false
	}
	digest := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(a.tokenHash[:], digest[:]) == 1
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok || !a.Verify(token) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="safeline-mcp"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
