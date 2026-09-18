package api

import (
	"encoding/json"
	"fmt"
	"io"
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

// maxVersionInfoBytes bounds the upgrade information this server reads from
// the platform. The body is parsed and logged on error, so a large or endless
// answer from a misbehaving (or hostile) endpoint must not be buffered whole.
const maxVersionInfoBytes = 1 << 20

type idsRequest struct {
	IDs []uint `json:"ids" form:"ids"`
}

type pageRequest struct {
	// PageSize is capped at 100: the value is handed to SQL LIMIT, so an
	// unbounded one lets a single request ask the server to load an arbitrary
	// number of rows. 100 is the largest page the console offers and the same
	// range the MCP tools accept.
	Page     int `json:"page"         form:"page,default=1"         binding:"min=1"`
	PageSize int `json:"page_size"    form:"page_size,default=10"   binding:"min=1,max=100"`
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
	// The response body has to be released on every path, otherwise the
	// connection (and its file descriptor) stays around until the transport
	// times it out.
	defer versionInfoRsp.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(versionInfoRsp.Body, maxVersionInfoBytes))
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
