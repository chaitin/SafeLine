package api

import (
	"net/http"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"chaitin.cn/dev/go/errors"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

type putSwitchRequest struct {
	IDs       []uint `json:"ids"`
	IsEnabled bool   `json:"is_enabled"`
}

func PostPolicyRule(ctx *gin.Context) {
	var params model.PolicyRule
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		policyRule := &model.PolicyRule{Action: params.Action, Comment: params.Comment, IsEnabled: params.IsEnabled, Pattern: params.Pattern}
		res := tx.Create(policyRule)
		if res.Error != nil {
			return res.Error
		}

		if err := fvm.PushFSL(tx); err != nil {
			return errors.New("Rules compile error, please check your params.")
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusInternalServerError)
		return
	}

	response.Success(ctx, nil)
}

func PutSwitchPolicyRule(ctx *gin.Context) {
	var params putSwitchRequest
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.PolicyRule{}).Where(params.IDs).Updates(model.PolicyRule{IsEnabled: params.IsEnabled})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("Data queried does not exist")
		}

		if err := fvm.PushFSL(tx); err != nil {
			return errors.New("Rules compile error, please check your params.")
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusInternalServerError)
		return
	}
}

func PutPolicyRule(ctx *gin.Context) {
	var params model.PolicyRule
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		var policyRule model.PolicyRule
		res := tx.Where(params.ID).First(&policyRule)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("Data queried does not exist")
		}

		policyRule.Action = params.Action
		policyRule.Comment = params.Comment
		policyRule.IsEnabled = params.IsEnabled
		policyRule.Pattern = params.Pattern
		tx.Save(&policyRule)

		if err := fvm.PushFSL(tx); err != nil {
			return errors.New("Rules compile error, please check your params.")
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusInternalServerError)
		return
	}

	response.Success(ctx, nil)
}

func DeletePolicyRule(ctx *gin.Context) {
	var params idsRequest
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where(params.IDs).Delete(&model.PolicyRule{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("Data queried does not exist")
		}

		if err := fvm.PushFSL(tx); err != nil {
			return errors.New("Rules compile error, please check your params.")
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusInternalServerError)
		return
	}

	response.Success(ctx, nil)
}

func GetPolicyRule(ctx *gin.Context) {
	var params pageRequest
	if err := ctx.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	var policyRuleList []model.PolicyRule
	db.Limit(params.PageSize).Offset(params.PageSize * (params.Page - 1)).Order("id desc").Find(&policyRuleList)

	var total int64
	db.Model(&model.PolicyRule{}).Count(&total)

	response.Success(ctx, gin.H{"data": policyRuleList, "total": total})
}
