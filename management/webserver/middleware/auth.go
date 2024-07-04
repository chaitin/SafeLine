package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get(constants.DefaultSessionUserKey)
	if user == nil {
		response.Error(c, response.ErrorLoginRequired, http.StatusUnauthorized)
		c.Abort()
		return
	}

	// extend session expired time
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: 3600 * 24 * 7,
		//Domain: options.Domain,
		//HttpOnly: true,
		//SameSite: http.SameSiteLaxMode,
		//Secure:   false,
	})

	if err := session.Save(); err != nil {
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when creating sessions"}, http.StatusInternalServerError)
		return
	}

	c.Next()
}

func ReadOnly(c *gin.Context) {
	if c.Request.Method != "GET" && c.Request.Method != "HEAD" && c.Request.Method != "OPTIONS" {
		response.Error(c, response.ErrorReadOnly, http.StatusBadRequest)
		c.Abort()
	}
}
