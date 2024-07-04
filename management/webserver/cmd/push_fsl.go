package cmd

import (
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm"
)

func PushFSL() error {
	return fvm.PushFSL(database.GetDB().DB)
}
