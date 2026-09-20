package utils

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const proxyName = "HTTPS_PROXY"

const (
	// httpClientTimeout bounds one request end to end, from the connection
	// attempt to the last byte of the body.
	httpClientTimeout = 15 * time.Second

	// httpClientDialTimeout bounds the connection attempt, and
	// httpClientTLSHandshakeTimeout the handshake that follows it.
	httpClientDialTimeout         = 5 * time.Second
	httpClientTLSHandshakeTimeout = 5 * time.Second

	// httpClientResponseHeaderTimeout bounds the wait for the response header
	// once the request was written.
	httpClientResponseHeaderTimeout = 10 * time.Second
)

var (
	httpClientOnce sync.Once
	httpClient     *http.Client
)

// GetHTTPClient returns the shared HTTP client.
//
// The client is built once: the API handlers and the cron jobs call this at the
// same time, and the unsynchronized lazy assignment used before let two of them
// build (and see) a client at the same time.
//
// Every phase of a request is bounded. The client is used by endpoints that any
// caller can reach, such as the upgrade tip lookup, and the endpoints it talks
// to are remote: an upstream that accepts the connection and then stops
// answering would otherwise hold the handler goroutine and its socket for as
// long as the connection lives, which a caller can repeat without limit.
func GetHTTPClient() *http.Client {
	httpClientOnce.Do(func() {
		// Certificates are verified: this client talks to the telemetry
		// endpoint and to the upgrade server, both of which are remote and can
		// be spoofed by a man in the middle when verification is disabled.
		tr := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   httpClientDialTimeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   httpClientTLSHandshakeTimeout,
			ResponseHeaderTimeout: httpClientResponseHeaderTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		}

		proxyUrl, existed := os.LookupEnv(proxyName)
		if existed {
			uri, _ := url.Parse(proxyUrl)
			tr.Proxy = http.ProxyURL(uri)
		}

		httpClient = &http.Client{Transport: tr, Timeout: httpClientTimeout}
	})

	return httpClient
}
