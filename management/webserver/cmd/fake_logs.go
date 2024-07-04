package cmd

import "chaitin.cn/patronus/safeline-2/management/webserver/model"

func FakeLogs() {
	model.InitDetectLogSamples()
}
