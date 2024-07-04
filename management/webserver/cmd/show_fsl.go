package cmd

import (
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm"
)

func ShowFSL() (string, error) {
	return fvm.GenerateFullFSL(database.GetDB().DB)
}
