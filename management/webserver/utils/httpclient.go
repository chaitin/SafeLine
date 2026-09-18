package utils

import (
	"net/http"
	"net/url"
	"os"
	"time"
)

const proxyName = "HTTPS_PROXY"

var httpClient *http.Client

func GetHTTPClient() *http.Client {
	if httpClient == nil {
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
	}

	return httpClient
}
