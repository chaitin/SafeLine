package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BootstrapTokenFileEnv names the file the bootstrap token is read from when
// BootstrapTokenEnv is empty.
const BootstrapTokenFileEnv = "MGT_BOOTSTRAP_TOKEN_FILE"

// DefaultBootstrapTokenFile is next to the control-channel token on the
// shared sock volume.
const DefaultBootstrapTokenFile = "/app/sock/tfa_bootstrap_token"

const bootstrapTokenBytes = 32
const bootstrapTokenFileMode os.FileMode = 0600

// bootstrapToken is the value GET /api/OTPUrl must present. It is never empty
// after a successful EnsureBootstrapToken.
var bootstrapToken string

// BootstrapTokenFile returns the path of the persisted bootstrap token.
func BootstrapTokenFile() string {
	if path := strings.TrimSpace(os.Getenv(BootstrapTokenFileEnv)); path != "" {
		return path
	}
	return DefaultBootstrapTokenFile
}

// EnsureBootstrapToken loads MGT_BOOTSTRAP_TOKEN or a persisted file, creating
// the file once when neither is set. The console must send the value as
// X-Bootstrap-Token; an empty token is never accepted.
func EnsureBootstrapToken() (created bool, path string, err error) {
	if token := strings.TrimSpace(os.Getenv(BootstrapTokenEnv)); token != "" {
		bootstrapToken = token
		return false, "", nil
	}

	path = BootstrapTokenFile()
	token, created, err := loadOrCreateBootstrapToken(path)
	if err != nil {
		return false, path, err
	}
	if token == "" {
		return false, path, fmt.Errorf("bootstrap token is empty")
	}
	bootstrapToken = token
	return created, path, nil
}

func loadOrCreateBootstrapToken(path string) (string, bool, error) {
	if content, err := os.ReadFile(path); err == nil {
		if token := strings.TrimSpace(string(content)); token != "" {
			return token, false, nil
		}
	}

	raw := make([]byte, bootstrapTokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", false, err
	}
	value := hex.EncodeToString(raw)

	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return "", false, err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, bootstrapTokenFileMode)
	if err != nil {
		if !os.IsExist(err) {
			return "", false, err
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return "", false, readErr
		}
		if token := strings.TrimSpace(string(content)); token != "" {
			return token, false, nil
		}
		// Empty leftover file: replace it. Do not truncate a non-empty secret.
		if err := os.WriteFile(path, []byte(value+"\n"), bootstrapTokenFileMode); err != nil {
			return "", false, err
		}
		return value, true, nil
	}
	if _, err = file.WriteString(value + "\n"); err != nil {
		_ = file.Close()
		return "", false, err
	}
	if err = file.Close(); err != nil {
		return "", false, err
	}
	return value, true, nil
}
