package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chaitin.cn/dev/go/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/rpc"
)

func publishWebsiteConfig(website *model.Website) error {
	byteWebsite, err := json.Marshal(website)
	if err != nil {
		return err
	}

	err = rpc.Publish(byteWebsite, rpc.EventTypeWebsite)
	if err != nil {
		return err
	}

	return nil
}

func publishDeleteWebsiteConfig(id []uint) error {
	byteId, err := json.Marshal(id)
	if err != nil {
		return err
	}

	err = rpc.Publish(byteId, rpc.EventTypeDeleteWebsite)
	if err != nil {
		return err
	}

	return nil
}

func PostWebsite(ctx *gin.Context) {
	var params model.Website
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		website := &model.Website{Comment: params.Comment, ServerNames: params.ServerNames, Upstreams: params.Upstreams, Ports: params.Ports,
			CertFilename: params.CertFilename, KeyFilename: params.KeyFilename, IsEnabled: true}
		res := tx.Create(website)
		if res.Error != nil {
			return res.Error
		}

		if config.GlobalConfig.Server.DevMode {
			return nil
		}

		err := publishWebsiteConfig(website)
		if err != nil {
			return err
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

func PutWebsite(ctx *gin.Context) {
	var params model.Website
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		var website model.Website
		res := tx.Where(params.ID).First(&website)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("Data queried does not exist")
		}

		website.Comment = params.Comment
		website.ServerNames = params.ServerNames
		website.Upstreams = params.Upstreams
		website.Ports = params.Ports
		website.CertFilename = params.CertFilename
		website.KeyFilename = params.KeyFilename
		website.IsEnabled = true
		tx.Save(&website)

		if config.GlobalConfig.Server.DevMode {
			return nil
		}

		err := publishWebsiteConfig(&params)
		if err != nil {
			return err
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

func DeleteWebsite(ctx *gin.Context) {
	var params idsRequest
	if err := ctx.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where(params.IDs).Delete(&model.Website{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("Data queried does not exist")
		}

		if config.GlobalConfig.Server.DevMode {
			return nil
		}

		err := publishDeleteWebsiteConfig(params.IDs)
		if err != nil {
			return err
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

func GetWebsite(ctx *gin.Context) {
	var params pageRequest
	if err := ctx.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(ctx, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	type Website struct {
		model.Website
		ReqValue    int64 `json:"req_value"`
		DeniedValue int64 `json:"denied_value"`
	}
	var websiteList []Website
	db.Limit(params.PageSize).Offset(params.PageSize * (params.Page - 1)).Order("id desc").Find(&websiteList)

	var statistics []model.SystemStatistics
	var websiteIds []string
	for _, i := range websiteList {
		websiteIds = append(websiteIds, strconv.Itoa(int(i.ID)))
	}
	db.Where("created_at >= date_trunc('day',now())").Where("website in (?)", websiteIds).Find(&statistics)

	for i, website := range websiteList {
		for _, j := range statistics {
			if strconv.Itoa(int(website.ID)) == j.Website {
				if j.Type == "website-req" {
					websiteList[i].ReqValue = j.Value
				} else if j.Type == "website-denied" {
					websiteList[i].DeniedValue = j.Value
				}
			}
		}
	}

	var total int64
	db.Model(&model.Website{}).Count(&total)

	response.Success(ctx, gin.H{"data": websiteList, "total": total})
}
