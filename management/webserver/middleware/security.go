package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets the response headers that keep a browser from
// interpreting an API answer as something else, and keep the console out of
// frames.
//
// The policy deliberately restricts only embedding, plugins and the base URI:
// the management console is served from the same host, so a stricter policy
// would have to be coordinated with the console bundle.
func SecurityHeaders(c *gin.Context) {
	// A JSON answer must never be sniffed into HTML or JavaScript.
	c.Header("X-Content-Type-Options", "nosniff")

	// The console must not be embedded, so that another page cannot click
	// through it on behalf of the administrator.
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "frame-ancestors 'none'; object-src 'none'; base-uri 'none'")

	// Keep the console URL out of the Referer header of third party requests.
	c.Header("Referrer-Policy", "no-referrer")

	// The console is reached over TLS; a browser that saw this header once will
	// refuse to fall back to plain HTTP for the following year. includeSubDomains
	// is deliberately left out: the console shares its host with whatever else
	// the deployment serves from it, and pinning every subdomain of that name is
	// not this service's decision to make.
	c.Header("Strict-Transport-Security", "max-age=31536000")

	c.Next()
}
