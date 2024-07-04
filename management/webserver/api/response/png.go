package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const StreamContentType = "application/octet-stream"

func PNG(c *gin.Context, bytes []byte) {
	c.Data(http.StatusOK, StreamContentType, bytes)
}
