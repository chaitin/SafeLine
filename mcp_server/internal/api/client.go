package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

const (
	// tokenHeader carries the deployment token of the configured instance.
	tokenHeader = "X-SLCE-API-TOKEN"

	maxErrorBodySize    = 4 << 10
	maxResponseBodySize = 4 << 20

	// maxRedirects bounds how many redirects one request may follow.
	maxRedirects = 10
)

// loadRootCAs reads the certificate authority the instance is issued by.
//
// It is a file rather than the system pool because the installation that the
// MCP server talks to usually carries the certificate it generated for itself,
// which no system pool knows about.
func loadRootCAs(path string) (*x509.CertPool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read ca_file failed")
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(content) {
		return nil, errors.New("ca_file does not contain a PEM encoded certificate")
	}

	return roots, nil
}

// Client is the transport used by a single SafeLine instance. Its request
// method is intentionally unexported; APIClient exposes only approved read-only
// operations.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	headers    http.Header
}

func newClient(baseURL string, timeout time.Duration, insecureSkipVerify bool, caFile, token string) (*Client, error) {
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
		// Default insecureSkipVerify is false (config.yaml). Self-signed
		// consoles should set ca_file to pin the instance CA. Turning this
		// on is an explicit operator choice and logs a warn below; it is
		// not the default skip-verify path.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify}, // #nosec G402 -- deployment-controlled compatibility setting
	}

	if caFile != "" {
		// The console of a SafeLine installation serves a certificate that no
		// public authority signed, so the choice used to be between switching
		// verification off and failing every request. Pinning the authority of
		// the instance keeps it on, which is what keeps the deployment token in
		// the request header out of the hands of anything on the path.
		roots, err := loadRootCAs(caFile)
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig.RootCAs = roots
	}

	if insecureSkipVerify {
		logger.With("base_url", parsedURL.Redacted()).Warn(
			"insecure_skip_verify is set: the certificate of this instance is not checked, " +
				"so anything on the path can read the API token of the deployment. " +
				"Set ca_file to the certificate authority of the instance instead.")
	}

	client := &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
			// Every request carries the deployment token in X-SLCE-API-TOKEN.
			// net/http drops the Authorization header when a redirect leaves
			// the origin, but it knows nothing about this header, so a 3xx
			// from the configured instance would hand the token to whatever
			// host the Location header names. Redirects are therefore followed
			// only while they stay on the configured origin.
			CheckRedirect: sameOriginRedirects(parsedURL),
		},
		headers: make(http.Header),
	}
	client.headers.Set("Accept", "application/json")
	client.headers.Set("User-Agent", "SafeLine-MCP/1.0")
	client.headers.Set(tokenHeader, token)
	return client, nil
}

// sameOriginRedirects returns the redirect policy of a client: a redirect is
// followed only while it stays on the scheme, host and port of the configured
// instance. Anything else fails the request instead of replaying the token.
func sameOriginRedirects(base *url.URL) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRedirects)
		}
		if !sameOrigin(base, req.URL) {
			return fmt.Errorf("refused to follow a redirect to %s: it would send the instance token to another origin", req.URL.Redacted())
		}

		return nil
	}
}

// sameOrigin reports whether target has the same scheme, host and port as base.
func sameOrigin(base, target *url.URL) bool {
	return strings.EqualFold(base.Scheme, target.Scheme) &&
		strings.EqualFold(base.Hostname(), target.Hostname()) &&
		effectivePort(base) == effectivePort(target)
}

// effectivePort returns the port a URL connects to, including the default port
// of its scheme.
func effectivePort(u *url.URL) string {
	if port := u.Port(); port != "" {
		return port
	}
	if strings.EqualFold(u.Scheme, "https") {
		return "443"
	}

	return "80"
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
