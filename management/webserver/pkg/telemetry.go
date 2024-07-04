package pkg

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

type WebsiteResult struct {
	Id        uint
	UpdatedAt time.Time
}

type TelemetryInfo struct {
	Id string `json:"id"`
}

type SafelineInfo struct {
	Id string `json:"id"`

	Type    string `json:"type"`
	Version string `json:"version"`

	ReqCnt         int  `json:"req_cnt"`
	DetectLogCnt   int  `json:"detect_log_cnt"`
	SiteCnt        int  `json:"site_cnt"`
	HealthySiteCnt int  `json:"healthy_site_cnt"`
	RuleCnt        int  `json:"rule_cnt"`
	IsHealthy      bool `json:"is_healthy"`

	Behavior    map[string]int  `json:"behavior"`
	Websites    []WebsiteResult `json:"websites"`
	WebsitesCnt int             `json:"websites_cnt"`
}

type TelemetryRequest struct {
	Telemetry TelemetryInfo `json:"telemetry"`
	Safeline  SafelineInfo  `json:"safeline"`
}

func GetUploadTimestamp() string {
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayNoon := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 12, 0, 0, 0, yesterday.Location())
	return strconv.FormatInt(yesterdayNoon.Unix(), 10)
}

func GetUploadNonce() string {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	randStr := utils.RandStr(6)
	return fmt.Sprintf("%s%s", now, randStr)
}

func DoPostTelemetry(client *http.Client, addr string, reader io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, "https://"+addr+constants.TelemetryEntryPoint, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", constants.ApplicationJson)
	req.Header.Set("Content-Type", constants.ApplicationJson)

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
