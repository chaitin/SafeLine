package utils

import (
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const proxyName = "HTTPS_PROXY"

var (
	httpClientOnce sync.Once
	httpClient     *http.Client
)

// GetHTTPClient returns the shared HTTP client.
//
// The client is built once: the API handlers and the cron jobs call this at the
// same time, and the unsynchronized lazy assignment used before let two of them
// build (and see) a client at the same time.
func GetHTTPClient() *http.Client {
	httpClientOnce.Do(func() {
		// Certificates are verified: this client talks to the telemetry
		// endpoint and to the upgrade server, both of which are remote and can
		// be spoofed by a man in the middle when verification is disabled.
		tr := &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		}

		proxyUrl, existed := os.LookupEnv(proxyName)
		if existed {
			uri, _ := url.Parse(proxyUrl)
			tr.Proxy = http.ProxyURL(uri)
		}

		httpClient = &http.Client{Transport: tr}
	})

	return httpClient
}
