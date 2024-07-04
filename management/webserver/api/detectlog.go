package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"chaitin.cn/dev/go/errors"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg"

	"github.com/gin-gonic/gin"

	"chaitin.cn/dev/go/log"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

type (
	GetDetectLogDetailRequest struct {
		EventId string `json:"event_id"   form:"event_id"`
	}

	PostFalsePositivesRequest struct {
		EventId string `json:"event_id"`
	}

	telemetryFalsePositives struct {
		Telemetry struct {
			Id string `json:"id"`
		} `json:"telemetry"`

		Safeline struct {
			Id        string          `json:"id"`
			Type      string          `json:"type"`
			DetectLog model.DetectLog `json:"detect_log"`
		} `json:"safeline"`
	}
)

func getDetectLog(eventId string) (*model.DetectLog, error) {
	db := database.GetDB()
	var detectLogBasic model.DetectLogBasic
	res := db.Where(&model.DetectLogBasic{EventId: eventId}).First(&detectLogBasic)
	if res.RowsAffected == 0 {
		return nil, errors.New("Data queried does not exist")
	}

	var detectLogDetail model.DetectLogDetail
	db.Where(&model.DetectLogDetail{EventId: eventId}).First(&detectLogDetail)

	detectLog, err := model.TransformDetectLog(&detectLogBasic, &detectLogDetail)
	if err != nil {
		return nil, err
	}

	return detectLog, nil
}

func GetDetectLogList(c *gin.Context) {
	var params pageRequest
	if err := c.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}
	db := database.GetDB()

	tx := db.Where("")
	// 按 ip 搜索条件
	if ip := c.Query("ip"); ip != "" {
		tx = tx.Where("src_ip = ?", ip)
	}
	// 按 url 搜索条件
	if url := c.Query("url"); url != "" {
		tx = tx.Where("url_path like ?", "%"+url+"%")
	}
	// 按 type 搜索条件
	if at := c.Query("attack_type"); at != "" {
		ns := make([]int, 0)
		for _, s := range strings.Split(at, ",") {
			n, err := strconv.Atoi(s)
			if err == nil {
				ns = append(ns, n)
			}
		}
		tx = tx.Where("attack_type in (?)", ns)
	}

	var total int64
	tx.Model(&model.DetectLogBasic{}).Count(&total)

	var basicList []model.DetectLogBasic
	tx.Limit(params.PageSize).Offset(params.PageSize * (params.Page - 1)).Order("id desc").Find(&basicList)
	var dLogList []*model.DetectLog
	for _, basic := range basicList {
		dLog, err := model.TransformDetectLog(&basic, nil)
		if err != nil {
			logger.Warn(err)
			continue
		}
		dLogList = append(dLogList, dLog)
	}
	response.Success(c, gin.H{"data": dLogList, "total": total})
}

func GetDetectLogDetail(c *gin.Context) {
	var params GetDetectLogDetailRequest
	if err := c.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	detectLog, err := getDetectLog(params.EventId)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorDataNotExist, http.StatusNotFound)
		return
	}

	response.Success(c, detectLog)
}

func PostFalsePositives(c *gin.Context) {
	var params PostFalsePositivesRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	detectLog, err := getDetectLog(params.EventId)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorDataNotExist, http.StatusNotFound)
		return
	}

	db := database.GetDB()
	var option model.Options
	db.Where(&model.Options{Key: constants.MachineID}).First(&option)

	var jsonData telemetryFalsePositives
	jsonData.Telemetry.Id = constants.TelemetryId
	jsonData.Safeline.Id = option.Value
	jsonData.Safeline.Type = constants.FalsePositives
	jsonData.Safeline.DetectLog = *detectLog

	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Warn(err)
		response.Success(c, nil)
		return
	}
	reader := bytes.NewReader(data)

	client := utils.GetHTTPClient()
	addr := config.GlobalConfig.Telemetry.Addr

	rsp, err := pkg.DoPostTelemetry(client, addr, reader)
	if err != nil {
		log.Warn(err)
		response.Success(c, nil)
		return
	}

	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		log.Warn(fmt.Sprintf("Transfer telemetry failed, status code = %d", rsp.StatusCode), err)
		response.Success(c, nil)
		return
	}

	response.Success(c, nil)
}
