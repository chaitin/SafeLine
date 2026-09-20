// Package sessionopts holds the cookie attributes of the management console
// session. The login handler and the authentication middleware share it so the
// flags cannot drift apart from each other.
package sessionopts

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
)

// maxAge is the lifetime of a management console session: one week.
const maxAge = 3600 * 24 * 7

// Options returns the cookie attributes used for the management console
// session.
//
// HttpOnly keeps the session cookie out of reach of JavaScript, SameSite=Lax
// prevents browsers from attaching the cookie to cross-site POST/PUT (which
// would otherwise turn the state-changing API into CSRF targets), and Secure
// keeps the cookie on HTTPS outside of a development setup.
//
// There is no separate CSRF header by design: the console is a cookie-session
// SPA, and Lax already drops the cookie on cross-site mutating requests.
// Adding a token would require a frontend change; do not treat the missing
// middleware as an accidental gap.
func Options(c *gin.Context) sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secureCookie(c),
		SameSite: http.SameSiteLaxMode,
	}
}

// secureCookie reports whether the session cookie may be restricted to HTTPS.
//
// The console is only reachable over TLS in production, but the webserver
// itself sits behind a TLS terminating nginx, so a plain HTTP request in front
// of it can still be a legitimate request of an HTTPS session.
func secureCookie(c *gin.Context) bool {
	if !config.GlobalConfig.Server.DevMode {
		return true
	}

	if c.Request.TLS != nil {
		return true
	}

	// Honour the scheme reported by a TLS terminating proxy, which is what the
	// development setup is expected to use when it wants to exercise the
	// production cookie attributes.
	proto := c.Request.Header.Get("X-Forwarded-Proto")
	if i := strings.IndexByte(proto, ','); i >= 0 {
		proto = proto[:i]
	}

	return strings.EqualFold(strings.TrimSpace(proto), "https")
}
