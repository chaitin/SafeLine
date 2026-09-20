package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// The client sends the deployment token in a header of its own, while net/http
// only strips the standard Authorization header when a redirect leaves the
// origin. A 3xx answered by the configured instance would therefore hand the
// token to whatever host the Location header names.
func TestClientRefusesCrossOriginRedirect(t *testing.T) {
	var tokenAtOtherOrigin atomic.Value
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenAtOtherOrigin.Store(r.Header.Get(tokenHeader))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"err":null,"data":{}}`))
	}))
	defer other.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/stolen", http.StatusFound)
	}))
	defer origin.Close()

	client, err := newClient(origin.URL, 5*time.Second, false, "deployment-token")
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	err = client.get(context.Background(), "/api/open/events", url.Values{}, &result)
	if err == nil {
		t.Fatal("a redirect to another origin has to fail the request")
	}
	if !strings.Contains(err.Error(), "refused to follow a redirect") {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tokenAtOtherOrigin.Load(); got != nil {
		t.Fatalf("the token reached another origin: %v", got)
	}
}

func TestClientFollowsSameOriginRedirect(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/open/events", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/moved", http.StatusFound)
	})
	var tokenAfterRedirect atomic.Value
	mux.HandleFunc("/moved", func(w http.ResponseWriter, r *http.Request) {
		tokenAfterRedirect.Store(r.Header.Get(tokenHeader))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"err":null,"data":{"ok":true}}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := newClient(server.URL, 5*time.Second, false, "deployment-token")
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := client.get(context.Background(), "/api/open/events", url.Values{}, &result); err != nil {
		t.Fatalf("a same-origin redirect has to be followed: %v", err)
	}
	if got := tokenAfterRedirect.Load(); got != "deployment-token" {
		t.Fatalf("the token was lost on the way: %v", got)
	}
}

func TestSameOrigin(t *testing.T) {
	cases := []struct {
		base, target string
		want         bool
	}{
		{"https://safeline.example:9443", "https://safeline.example:9443/other", true},
		{"https://safeline.example:9443", "https://SAFELINE.example:9443/other", true},
		{"https://safeline.example", "https://safeline.example:443/other", true},
		{"http://safeline.example", "http://safeline.example:80/other", true},
		{"https://safeline.example", "http://safeline.example/other", false},
		{"https://safeline.example", "https://safeline.example:8443/other", false},
		{"https://safeline.example:9443", "https://evil.example:9443/other", false},
		{"https://safeline.example:9443", "https://safeline.example.evil.example:9443/other", false},
	}

	for _, c := range cases {
		base, err := url.Parse(c.base)
		if err != nil {
			t.Fatal(err)
		}
		target, err := url.Parse(c.target)
		if err != nil {
			t.Fatal(err)
		}
		if got := sameOrigin(base, target); got != c.want {
			t.Errorf("sameOrigin(%s, %s) = %v, want %v", c.base, c.target, got, c.want)
		}
	}
}

// TestPlainClientLeaksTokenAcrossRedirects documents what the policy above
// prevents: with the default redirect handling of net/http the custom token
// header arrives at the other origin.
func TestPlainClientLeaksTokenAcrossRedirects(t *testing.T) {
	var tokenAtOtherOrigin atomic.Value
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenAtOtherOrigin.Store(r.Header.Get(tokenHeader))
		w.Write([]byte(`{}`))
	}))
	defer other.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/stolen", http.StatusFound)
	}))
	defer origin.Close()

	req, err := http.NewRequest(http.MethodGet, origin.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(tokenHeader, "deployment-token")
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got := tokenAtOtherOrigin.Load(); got != "deployment-token" {
		t.Fatalf("the plain client was expected to leak the token, got %v", got)
	}
}
