package cron

import (
	"github.com/robfig/cron/v3"

	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/log"
)

var logger = log.GetLogger("cron")

func newCronWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func StartCron() error {
	cronInstance := newCronWithSeconds()

	_, err := cronInstance.AddFunc(specCheckForbiddenPage, checkAndUpdateForbiddenPage)
	if err != nil {
		return err
	}

	cronInstance.Start()
	return nil
}
