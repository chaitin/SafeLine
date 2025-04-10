package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

// Client API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// ClientOption Client configuration options
type ClientOption func(*Client)

// WithTimeout Set timeout duration
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHeader Set request header
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithBaseURL Set base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithInsecureSkipVerify Set whether to skip certificate verification
func WithInsecureSkipVerify(skip bool) ClientOption {
	return func(c *Client) {
		if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			if transport.TLSClientConfig == nil {
				transport.TLSClientConfig = &tls.Config{}
			}
			transport.TLSClientConfig.InsecureSkipVerify = skip
		}
	}
}

// NewClient Create new API client
func NewClient(opts ...ClientOption) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	c := &Client{
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Request Send request
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "marshal request body failed")
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}
	logger.With("url", reqURL).Debug("request url")
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}

	// Set common headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request failed")
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response body failed")
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	// Parse response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err == nil {
			return nil
		}
		var respData map[string]interface{}
		if err := json.Unmarshal(respBody, &respData); err != nil {
			return errors.Wrap(err, "unmarshal response failed")
		}
		if respData["err"] != nil || respData["msg"] != nil {
			return errors.New(respData["msg"].(string))
		}
	}

	return nil
}

// Get Send GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, nil, result)
}

// Post Send POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, body, result)
}

// Put Send PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, path, body, result)
}

// Delete Send DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, nil, result)
}
