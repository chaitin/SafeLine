package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOriginMatchesRequestIgnoresScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "http://127.0.0.1:9001/api/Login", nil)
	c.Request.Host = "example.com:9443"

	origin, err := url.Parse("https://example.com:9443")
	if err != nil {
		t.Fatal(err)
	}
	if !originMatchesRequest(origin, c) {
		t.Fatal("https Origin should match Host behind HTTP gin")
	}

	evil, err := url.Parse("https://evil.example")
	if err != nil {
		t.Fatal(err)
	}
	if originMatchesRequest(evil, c) {
		t.Fatal("foreign Origin must not match")
	}
}
