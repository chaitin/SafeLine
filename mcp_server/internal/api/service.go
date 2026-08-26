package api

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chaitin/SafeLine/mcp_server/internal/config"
	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

const attackEventsPath = "/api/open/events"
const maxTokenFileSize = 16 << 10

// InstanceSource is safe to include in tool output and identifies the exact
// SafeLine instance from which the data was read.
type InstanceSource struct {
	InstanceID  string `json:"instance_id"`
	DisplayName string `json:"display_name"`
}

// APIClient is an independently authenticated client for one SafeLine instance.
type APIClient struct {
	client *Client
	source InstanceSource
}

// ClientRegistry maps stable instance IDs to independently configured clients.
// Credentials are loaded once per instance and are never shared across entries.
type ClientRegistry struct {
	clients map[string]*APIClient
	order   []string
}

// NewClientRegistry constructs a registry without mutating global state.
func NewClientRegistry(configs []*config.InstanceConfig) (*ClientRegistry, error) {
	if len(configs) == 0 {
		return nil, errors.New("at least one SafeLine instance is required")
	}

	registry := &ClientRegistry{
		clients: make(map[string]*APIClient, len(configs)),
		order:   make([]string, 0, len(configs)),
	}
	for index, cfg := range configs {
		if cfg == nil {
			return nil, errors.New(fmt.Sprintf("instances[%d] is required", index))
		}
		if strings.TrimSpace(cfg.ID) == "" {
			return nil, errors.New(fmt.Sprintf("instances[%d].id is required", index))
		}
		if _, exists := registry.clients[cfg.ID]; exists {
			return nil, errors.New(fmt.Sprintf("duplicate instance id %q", cfg.ID))
		}

		client, err := newAPIClient(cfg)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("initialize instance %q failed", cfg.ID))
		}
		registry.clients[cfg.ID] = client
		registry.order = append(registry.order, cfg.ID)
	}
	return registry, nil
}

// Get returns the client for instanceID without falling back to another
// instance. This prevents a typo from accidentally querying with the wrong
// instance credential.
func (r *ClientRegistry) Get(instanceID string) (*APIClient, error) {
	if r == nil {
		return nil, errors.New("API service not initialized")
	}
	client, ok := r.clients[instanceID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("unknown SafeLine instance %q", instanceID))
	}
	return client, nil
}

// Sources lists configured instances without exposing addresses or credentials.
func (r *ClientRegistry) Sources() []InstanceSource {
	if r == nil {
		return nil
	}
	sources := make([]InstanceSource, 0, len(r.order))
	for _, id := range r.order {
		sources = append(sources, r.clients[id].source)
	}
	return sources
}

var (
	registryMu     sync.RWMutex
	activeRegistry *ClientRegistry
)

// InitInstances initializes the runtime registry used by MCP tools.
func InitInstances(configs []*config.InstanceConfig) error {
	registry, err := NewClientRegistry(configs)
	if err != nil {
		logger.With("error", err).Error("failed to initialize API service")
		return err
	}

	registryMu.Lock()
	activeRegistry = registry
	registryMu.Unlock()
	logger.With("instances", len(configs)).Info("API service initialized successfully")
	return nil
}

// ClientFor returns the explicitly selected instance client.
func ClientFor(instanceID string) (*APIClient, error) {
	registryMu.RLock()
	registry := activeRegistry
	registryMu.RUnlock()
	return registry.Get(instanceID)
}

func newAPIClient(cfg *config.InstanceConfig) (*APIClient, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, errors.New("base_url is required")
	}

	token, err := loadToken(cfg)
	if err != nil {
		return nil, err
	}
	timeout := 30
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}
	client, err := newClient(cfg.BaseURL, time.Duration(timeout)*time.Second, cfg.InsecureSkipVerify, token)
	if err != nil {
		return nil, err
	}

	displayName := strings.TrimSpace(cfg.DisplayName)
	if displayName == "" {
		displayName = cfg.ID
	}
	return &APIClient{
		client: client,
		source: InstanceSource{InstanceID: cfg.ID, DisplayName: displayName},
	}, nil
}

func loadToken(cfg *config.InstanceConfig) (string, error) {
	if cfg.TokenFile != "" {
		dirInfo, err := os.Lstat(filepath.Dir(cfg.TokenFile))
		if err != nil {
			return "", errors.Wrap(err, "inspect token_file directory failed")
		}
		if dirInfo.Mode()&os.ModeSymlink != 0 || !dirInfo.IsDir() || dirInfo.Mode().Perm()&0o022 != 0 {
			return "", errors.New("token_file directory must be a real directory that is not group or world writable")
		}

		pathInfo, err := os.Lstat(cfg.TokenFile)
		if err != nil {
			return "", errors.Wrap(err, "inspect token_file failed")
		}
		if pathInfo.Mode()&os.ModeSymlink != 0 || !pathInfo.Mode().IsRegular() || pathInfo.Mode().Perm()&0o077 != 0 {
			return "", errors.New("token_file must be a regular file with no group or other permissions")
		}

		file, err := os.Open(cfg.TokenFile)
		if err != nil {
			return "", errors.Wrap(err, "open token_file failed")
		}
		defer file.Close()
		openedInfo, err := file.Stat()
		if err != nil {
			return "", errors.Wrap(err, "inspect opened token_file failed")
		}
		if !os.SameFile(pathInfo, openedInfo) || !openedInfo.Mode().IsRegular() || openedInfo.Mode().Perm()&0o077 != 0 {
			return "", errors.New("token_file changed while opening")
		}

		contents, err := io.ReadAll(io.LimitReader(file, maxTokenFileSize+1))
		if err != nil {
			return "", errors.Wrap(err, "read token_file failed")
		}
		if len(contents) > maxTokenFileSize {
			return "", errors.New("token_file exceeds the maximum size")
		}
		token := strings.TrimSpace(string(contents))
		if token == "" {
			return "", errors.New("token_file is empty")
		}
		return token, nil
	}

	return "", errors.New("token_file is required")
}

// Source returns the non-sensitive identity included in tool results.
func (c *APIClient) Source() InstanceSource {
	return c.source
}

// GetAttackEvents is the only SafeLine management request exposed by the
// client. The path and HTTP method are fixed rather than supplied by a tool.
func (c *APIClient) GetAttackEvents(ctx context.Context, query url.Values, result any) error {
	return c.client.get(ctx, attackEventsPath, query, result)
}
