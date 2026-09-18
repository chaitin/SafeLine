package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm"
)

func PutPolicyGroupGlobal(ctx *gin.Context) {
	var params model.PolicyGroup
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	// Reject unknown modules and modes before anything is persisted: the map
	// keys are rendered into the skynet configuration, so a bad key has to be
	// reported to the caller instead of reaching the engine or the renderer.
	if err := model.ValidatePolicyGroup(params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrorParamNotOK.Err, Msg: err.Error()}, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		var pggOption model.Options
		res := tx.Where(&model.Options{Key: constants.PolicyGroupGlobal}).First(&pggOption)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errDataNotExist
		}

		pggStr, err := json.Marshal(params)
		if err != nil {
			logger.Error(err)
			response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: responseMessage(err)}, http.StatusInternalServerError)
			return err
		}

		pggOption.Value = string(pggStr)
		tx.Save(&pggOption)

		if err := fvm.PushFSL(tx); err != nil {
			return errRulesCompile
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: responseMessage(err)}, http.StatusInternalServerError)
		return
	}

	response.Success(ctx, nil)
}

func GetPolicyGroupGlobal(ctx *gin.Context) {
	var pggOption model.Options
	database.GetDB().Where(&model.Options{Key: constants.PolicyGroupGlobal}).First(&pggOption)

	var pgg model.PolicyGroup
	err := json.Unmarshal([]byte(pggOption.Value), &pgg)
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: responseMessage(err)}, http.StatusInternalServerError)
		return
	}

	response.Success(ctx, gin.H{"data": pgg})
}

func PutSrcIPConfig(ctx *gin.Context) {
	var params model.SrcIPConfig
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		var scOption model.Options
		res := tx.Where(&model.Options{Key: constants.SrcIPConfig}).First(&scOption)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errDataNotExist
		}

		scStr, err := json.Marshal(params)
		if err != nil {
			logger.Error(err)
			response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: responseMessage(err)}, http.StatusInternalServerError)
			return err
		}

		scOption.Value = string(scStr)
		tx.Save(&scOption)

		if err := fvm.PushFSL(tx); err != nil {
			return errRulesCompile
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: responseMessage(err)}, http.StatusInternalServerError)
		return
	}

	response.Success(ctx, nil)
}

func GetSrcIPConfig(ctx *gin.Context) {
	srcIPConfig, err := model.GetSrcIPConfig(database.GetDB().DB)
	if err != nil {
		return
	}
	response.Success(ctx, gin.H{"data": srcIPConfig})
}
