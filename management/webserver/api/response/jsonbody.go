package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONBody struct {
	Err  string
	Msg  string
	Data interface{}
}

var (
	ErrorLoginRequired = JSONBody{ErrLoginRequired, "Login required", nil}
	ErrorParamNotOK    = JSONBody{ErrInternalError, "Error occurred when extracting params", nil}
	ErrorDataNotExist  = JSONBody{ErrDataNotExist, "Data queried does not exist", nil}
	ErrorReadOnly      = JSONBody{ErrReadOnly, "This environment is read only", nil}

	// ErrorInternal is the answer to a failed database or runtime call. The
	// text of the underlying error stays in the server log: it can carry
	// schema names, addresses and other details of the deployment that the
	// caller has no use for.
	ErrorInternal = JSONBody{ErrInternalError, "Error occurred when processing the request, please check the server log", nil}
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"msg":  "",
		"err":  nil,
	})
	//c.Abort()
}

func SuccessWithList(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"msg":  "",
		"err":  nil,
	})
	//c.Abort()
}

func SuccessWithMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"msg":  msg,
		"err":  nil,
	})
	//c.Abort()
}

func Error(c *gin.Context, rsp JSONBody, status int) {
	c.JSON(status, gin.H{
		"data": rsp.Data,
		"msg":  rsp.Msg,
		"err":  rsp.Err,
	})
	//c.Abort()
}
