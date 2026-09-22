package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
)

// SameOrigin rejects cross-site state-changing requests using the Origin
// header. The management console is a cookie-session SPA; Lax already drops
// the cookie on cross-site POST, and this check covers non-browser clients
// that still send Origin, plus browsers that attach the cookie on same-site
// confusing cases.
//
// Safe methods are left alone. Requests without a usable Origin are refused:
// Referrer-Policy is no-referrer, so Referer is not a fallback.
func SameOrigin(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		c.Next()
		return
	}

	origin := strings.TrimSpace(c.GetHeader("Origin"))
	if origin == "" || strings.EqualFold(origin, "null") {
		response.Error(c, response.JSONBody{Err: response.ErrForbidden, Msg: "cross-origin request is not allowed"}, http.StatusForbidden)
		c.Abort()
		return
	}

	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		response.Error(c, response.JSONBody{Err: response.ErrForbidden, Msg: "cross-origin request is not allowed"}, http.StatusForbidden)
		c.Abort()
		return
	}

	if !originMatchesRequest(parsed, c) {
		response.Error(c, response.JSONBody{Err: response.ErrForbidden, Msg: "cross-origin request is not allowed"}, http.StatusForbidden)
		c.Abort()
		return
	}

	c.Next()
}

func originMatchesRequest(origin *url.URL, c *gin.Context) bool {
	// Scheme is not compared: the console is reached over HTTPS while gin
	// often sees HTTP behind the container TLS terminator. Host still has
	// to match, so a foreign site cannot satisfy this check.
	if !strings.EqualFold(origin.Scheme, "http") && !strings.EqualFold(origin.Scheme, "https") {
		return false
	}
	return strings.EqualFold(origin.Host, c.Request.Host)
}
