package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rogpeppe/go-internal/semver"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

const VersionInfoEntrypoint = "/release/latest/version.json"

type idsRequest struct {
	IDs []uint `json:"ids" form:"ids"`
}

type pageRequest struct {
	Page     int `json:"page"         form:"page,default=1"         binding:"min=1"`
	PageSize int `json:"page_size"    form:"page_size,default=10"   binding:"min=1"`
}

type versionInfoResponse struct {
	LatestVersion string `json:"latest_version"`
	RecVersion    string `json:"rec_version"`
}

func GetVersion(c *gin.Context) {
	response.Success(c, gin.H{"version": strings.TrimPrefix(constants.Version, "ce-")})
}

func GetUpgradeTips(ctx *gin.Context) {
	client := utils.GetHTTPClient()
	logger.Debugf("GetUpgradeTips: %s", config.GlobalConfig.PlatformAddr+VersionInfoEntrypoint)
	versionInfoReq, err := http.NewRequest(http.MethodGet, config.GlobalConfig.PlatformAddr+VersionInfoEntrypoint, nil)
	if err != nil {
		logger.Warn(err)
		response.Success(ctx, gin.H{"upgrade_tips": constants.NotUpgrade})
		return
	}

	versionInfoRsp, err := client.Do(versionInfoReq)
	if err != nil {
		logger.Warn(err)
		response.Success(ctx, gin.H{"upgrade_tips": constants.NotUpgrade})
		return
	}
	body, err := ioutil.ReadAll(versionInfoRsp.Body)
	if err != nil {
		logger.Warn(err)
		response.Success(ctx, gin.H{"upgrade_tips": constants.NotUpgrade})
		return
	}

	versionInfo := &versionInfoResponse{}
	err = json.Unmarshal(body, versionInfo)
	if err != nil {
		logger.Warnf("err: %v, body: %s", err, body)
		response.Success(ctx, gin.H{"upgrade_tips": constants.NotUpgrade})
		return
	}

	currentVersion := fmt.Sprintf("v%s", constants.Version)
	latestVersionCmp := semver.Compare(currentVersion, versionInfo.LatestVersion)
	recVersionCmp := semver.Compare(currentVersion, versionInfo.RecVersion)
	if semver.Compare(versionInfo.LatestVersion, versionInfo.RecVersion) == -1 || latestVersionCmp == 1 {
		logger.Warnf("The version number is invalid, current version: %s, latest version: %s, rec version: %s",
			currentVersion, versionInfo.LatestVersion, versionInfo.RecVersion)
		response.Success(ctx, gin.H{"upgrade_tips": constants.NotUpgrade})
		return
	}

	var upgradeTips int
	if recVersionCmp == -1 {
		upgradeTips = constants.MustUpgrade
	} else if recVersionCmp == 0 {
		if latestVersionCmp == 0 {
			upgradeTips = constants.NotUpgrade
		} else {
			upgradeTips = constants.RecommendedUpgrade
		}
	} else {
		if latestVersionCmp < 0 {
			upgradeTips = constants.RecommendedUpgrade
		} else {
			upgradeTips = constants.NotUpgrade
		}
	}

	response.Success(ctx, gin.H{"upgrade_tips": upgradeTips})
}
