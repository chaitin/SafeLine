package cron

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"chaitin.cn/patronus/safeline-2/management/webserver/utils"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
)

// SpecUpdatePolicy http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html
// Seconds Minutes Hours Day-of-Month Month Day-of-Week Year (optional field)
// every 5 second (starting from 0s)
const SpecUpdatePolicy = "0/5 * * * * ?"

type statResponseBody struct {
	PolicyVersion string `json:"policy_version"`
}

func getBytecode() ([]byte, error) {
	buff, err := ioutil.ReadFile(config.GlobalConfig.Detector.FslBytecode)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func CheckAndUpdatePolicy() {
	if config.GlobalConfig.Detector.Addr == "" {
		return
	}

	exist, err := utils.FileExist(config.GlobalConfig.Detector.FslBytecode)
	if !exist || err != nil {
		return
	}

	data, err := getBytecode()
	if err != nil {
		logger.Error(err)
		return
	}

	reader := bytes.NewReader(data)

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}

	addr := config.GlobalConfig.Detector.Addr
	statReq, err := http.NewRequest(http.MethodGet, addr+constants.StatEntrypoint, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	updateReq, err := http.NewRequest(http.MethodPost, addr+constants.UpdateEntrypoint, reader)
	if err != nil {
		logger.Error(err)
		return
	}

	updateReq.Header.Set("Content-Type", constants.ContentType)

	statRspData := statResponseBody{}
	statRsp, err := client.Do(statReq)
	if err != nil || statRsp.StatusCode != http.StatusOK {
		logger.Warn(err)
		return
	}

	body, err := ioutil.ReadAll(statRsp.Body)
	if err != nil {
		logger.Warn(err)
		return
	}

	err = json.Unmarshal(body, &statRspData)
	if err != nil {
		logger.Warn(err)
		return
	}

	if statRspData.PolicyVersion == constants.DefaultPolicyVersion {
		return
	}

	err = statRsp.Body.Close()
	if err != nil {
		logger.Warn(err)
	}

	logger.Info("Update fsl bytecode")
	updateRsp, err := client.Do(updateReq)
	if err != nil {
		logger.Warn(err)
		return
	}

	if updateRsp.StatusCode != http.StatusOK {
		logger.Warnf("%s update policy, return %d", addr, updateRsp.StatusCode)
	}
	err = updateRsp.Body.Close()
	if err != nil {
		logger.Warn(err)
	}
}
