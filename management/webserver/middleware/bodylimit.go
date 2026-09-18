package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxRequestBodyBytes bounds the size of a request body the management API
// accepts.
//
// This server terminates TLS itself and its port is published on the host, so
// no proxy caps the body in front of it: without a limit a client can make it
// buffer an arbitrary amount of data. The largest body the API takes is an SSL
// certificate upload, which is a few tens of kilobytes, so the limit leaves a
// wide margin while still bounding what a single request may consume.
const MaxRequestBodyBytes = 16 << 20 // 16 MiB

// BodyLimit caps the request body at maxBytes.
//
// The body is replaced by a reader that reports an error once the limit is
// exceeded, so a handler that parses the body fails with a parameter error
// instead of reading it all. Requests without a body are not affected.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}
