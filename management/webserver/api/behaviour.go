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

func PostBehaviour(c *gin.Context) {
	var params PostBehaviourRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}
	db := database.GetDB()
	db.Create(&model.Behaviour{SrcRouter: params.SrcRouter, DstRouter: params.DstRouter})
	response.Success(c, nil)
}
