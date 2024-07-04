package cmd

import (
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

func ResetUser(username string) {
	db := database.GetDB()
	var user model.User
	db.Where(&model.User{Username: username}).First(&user)
	user.LastLoginTime = 0
	db.Save(&user)
}
