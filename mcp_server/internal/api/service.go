package api

import (
	"context"
	"sync"
	"time"

	"github.com/chaitin/SafeLine/mcp_server/internal/config"
	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

// APIClient API client implementation
type APIClient struct {
	client *Client
	config *config.APIConfig
}

var (
	instance *APIClient
	once     sync.Once
)

// Init Initialize API service
func Init(cfg *config.APIConfig) error {
	var err error
	once.Do(func() {
		instance, err = newAPIClient(cfg)
		if err != nil {
			logger.With("error", err).Error("failed to initialize API service")
			return
		}
		logger.Info("API service initialized successfully")
	})
	return err
}

// Service Get API service instance
func Service() *APIClient {
	if instance == nil {
		logger.Error("API service not initialized")
		panic("API service not initialized")
	}
	return instance
}

// newAPIClient Create new API client
func newAPIClient(config *config.APIConfig) (*APIClient, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	if config.BaseURL == "" {
		return nil, errors.New("base_url is required")
	}

	timeout := 30
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	opts := []ClientOption{
		WithBaseURL(config.BaseURL),
		WithTimeout(time.Duration(timeout) * time.Second),
		WithHeader("User-Agent", "SafeLine-MCP/1.0"),
		WithInsecureSkipVerify(config.InsecureSkipVerify),
	}

	// If token is configured, add authentication header
	if config.Token != "" {
		opts = append(opts, WithHeader("X-SLCE-API-TOKEN", config.Token))
	}

	client := NewClient(opts...)

	return &APIClient{
		client: client,
		config: config,
	}, nil
}

// Post Send POST request
func (c *APIClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.client.Request(ctx, "POST", path, body, result)
}

// Get Send GET request
func (c *APIClient) Get(ctx context.Context, path string, result interface{}) error {
	return c.client.Request(ctx, "GET", path, nil, result)
}

// Put Send PUT request
func (c *APIClient) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.client.Request(ctx, "PUT", path, body, result)
}

// Delete Send DELETE request
func (c *APIClient) Delete(ctx context.Context, path string, result interface{}) error {
	return c.client.Request(ctx, "DELETE", path, nil, result)
}
