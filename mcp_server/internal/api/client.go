package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

const (
	maxErrorBodySize    = 4 << 10
	maxResponseBodySize = 4 << 20
)

// Client is the transport used by a single SafeLine instance. Its request
// method is intentionally unexported; APIClient exposes only approved read-only
// operations.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	headers    http.Header
}

func newClient(baseURL string, timeout time.Duration, insecureSkipVerify bool, token string) (*Client, error) {
	parsedURL, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return nil, errors.Wrap(err, "parse base_url failed")
	}
	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return nil, errors.New("base_url must be an absolute http or https URL")
	}
	if parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return nil, errors.New("base_url must not contain a query or fragment")
	}
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify}, // #nosec G402 -- deployment-controlled compatibility setting
	}
	client := &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		headers: make(http.Header),
	}
	client.headers.Set("Accept", "application/json")
	client.headers.Set("User-Agent", "SafeLine-MCP/1.0")
	client.headers.Set("X-SLCE-API-TOKEN", token)
	return client, nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, result any) error {
	endpoint := *c.baseURL
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + path
	endpoint.RawPath = ""
	endpoint.RawQuery = query.Encode()

	logger.With("url", endpoint.String()).Debug("request url")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}
	req.Header = c.headers.Clone()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodySize+1))
	if err != nil {
		return errors.Wrap(err, "read response body failed")
	}
	if len(respBody) > maxResponseBodySize {
		return errors.New("SafeLine response body exceeds the 4 MiB limit")
	}

	var envelope struct {
		Err any    `json:"err"`
		Msg string `json:"msg"`
	}
	envelopeDecoded := json.Unmarshal(respBody, &envelope) == nil
	if envelopeDecoded {
		if err := responseError(envelope.Err, envelope.Msg); err != nil {
			return err
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body := strings.TrimSpace(string(respBody))
		if len(body) > maxErrorBodySize {
			body = body[:maxErrorBodySize] + "..."
		}
		return errors.New(fmt.Sprintf("SafeLine request failed with status %d: %s", resp.StatusCode, body))
	}

	if result == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, result); err != nil {
		return errors.Wrap(err, "unmarshal response failed")
	}
	return nil
}
