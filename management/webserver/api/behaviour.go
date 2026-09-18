package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

type PostBehaviourRequest struct {
	model.Behaviour
}

// maxRouterLength bounds the console routes that are stored. The values come
// from the browser, so they are limited to what the console can produce.
const maxRouterLength = 256

func PostBehaviour(c *gin.Context) {
	var params PostBehaviourRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	if len(params.SrcRouter) > maxRouterLength || len(params.DstRouter) > maxRouterLength {
		logger.Warnf("Rejected behaviour entry with %d/%d bytes of route", len(params.SrcRouter), len(params.DstRouter))
		response.Error(c, response.ErrorParamNotOK, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if err := db.Create(&model.Behaviour{SrcRouter: params.SrcRouter, DstRouter: params.DstRouter}).Error; err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorInternal, http.StatusInternalServerError)
		return
	}

	response.Success(c, nil)
}
