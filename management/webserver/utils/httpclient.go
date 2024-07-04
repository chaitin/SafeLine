package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"time"
)

const proxyName = "HTTPS_PROXY"

var httpClient *http.Client

func GetHTTPClient() *http.Client {
	if httpClient == nil {
		tr := &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
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
