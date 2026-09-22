package api

import (
	"encoding/json"
	"net/http"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

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

	if err := validatePolicyRulePatterns(params.Pattern); err != nil {
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusBadRequest)
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

func PutSwitchPolicyRule(ctx *gin.Context) {
	var params putSwitchRequest
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		// IsEnabled=false is the Go zero value. GORM skips zero-value fields in
		// Updates(struct), so a disable request would leave the row enabled.
		res := tx.Model(&model.PolicyRule{}).Where(params.IDs).
			Select("is_enabled").
			Updates(model.PolicyRule{IsEnabled: params.IsEnabled})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errDataNotExist
		}

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
}

func PutPolicyRule(ctx *gin.Context) {
	var params model.PolicyRule
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	if err := validatePolicyRulePatterns(params.Pattern); err != nil {
		response.Error(ctx, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusBadRequest)
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
			return errDataNotExist
		}

		policyRule.Action = params.Action
		policyRule.Comment = params.Comment
		policyRule.IsEnabled = params.IsEnabled
		policyRule.Pattern = params.Pattern
		tx.Save(&policyRule)

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
			return errDataNotExist
		}

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

// validatePolicyRulePatterns rejects a rule that would compile to a selector
// with no WHERE. An empty pattern list matches every request.
func validatePolicyRulePatterns(raw datatypes.JSON) error {
	if len(raw) == 0 {
		return errPolicyRuleEmpty
	}
	var patterns []model.PolicyRulePattern
	if err := json.Unmarshal(raw, &patterns); err != nil || len(patterns) == 0 {
		return errPolicyRuleEmpty
	}
	return nil
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
