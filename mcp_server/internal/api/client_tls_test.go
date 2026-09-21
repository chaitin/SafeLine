package api

import (
	"context"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The console of a SafeLine installation serves the certificate it generated
// for itself, so an MCP deployment that wants verification on has to be able to
// name the authority of that instance. ca_file is that, and it exists so that
// insecure_skip_verify is not the only way to reach an installation with a
// self-signed certificate.
func TestClientUsesTheAuthorityOfTheInstance(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"err":null,"data":{"ok":true}}`))
	}))
	defer server.Close()

	// The certificate is not known without the authority, and the request is
	// refused rather than sent in the clear or unverified.
	untrusted, err := newClient(server.URL, 5*time.Second, "", "deployment-token")
	if err != nil {
		t.Fatal(err)
	}

	var result struct {
		Data struct {
			OK bool `json:"ok"`
		} `json:"data"`
		Err any `json:"err"`
	}
	if err := untrusted.get(context.Background(), "/api/open/events", url.Values{}, &result); err == nil {
		t.Fatal("the instance was reached without trusting its certificate authority")
	}

	caPath := filepath.Join(t.TempDir(), "instance.crt")
	encoded := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	if err := os.WriteFile(caPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	trusted, err := newClient(server.URL, 5*time.Second, caPath, "deployment-token")
	if err != nil {
		t.Fatal(err)
	}

	if err := trusted.get(context.Background(), "/api/open/events", url.Values{}, &result); err != nil {
		t.Fatalf("the request with the authority of the instance failed: %v", err)
	}
	if !result.Data.OK || result.Err != nil {
		t.Fatalf("the request did not reach the instance: %+v", result)
	}
}

func TestClientRefusesACAFileWithoutACertificate(t *testing.T) {
	caPath := filepath.Join(t.TempDir(), "not-a-certificate.crt")
	if err := os.WriteFile(caPath, []byte("this is not a certificate\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := newClient("https://instance.example:9443", 5*time.Second, caPath, "token"); err == nil {
		t.Fatal("a ca_file without a certificate was accepted")
	}
}

func TestClientReportsAMissingCAFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.crt")
	if _, err := newClient("https://instance.example:9443", 5*time.Second, missing, "token"); err == nil {
		t.Fatal("a missing ca_file was accepted")
	}
}

func TestClientRejectsNonLoopbackHTTP(t *testing.T) {
	_, err := newClient("http://example.com:9443", 5*time.Second, "", "token")
	if err == nil || err.Error() != "deployment token must not be sent over non-loopback HTTP" {
		t.Fatalf("error = %v", err)
	}
}

func TestClientAllowsLoopbackHTTP(t *testing.T) {
	if _, err := newClient("http://127.0.0.1:9443", 5*time.Second, "", "token"); err != nil {
		t.Fatal(err)
	}
	if _, err := newClient("http://localhost:9443", 5*time.Second, "", "token"); err != nil {
		t.Fatal(err)
	}
}
