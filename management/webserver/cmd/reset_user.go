package cmd

import (
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/log"
)

var logger = log.GetLogger("cmd")

// ResetUser clears the credentials of a user so that the next visit of the
// console binds a new TFA secret.
//
// The pending secret and the address it was issued to are cleared together:
// a secret that was handed to somebody else cannot be bound by the
// administrator, and starting the bootstrap over from the shell is the only
// way out of that state.
//
// It also opens the anonymous bootstrap window again, because this command is
// how an operator asks for it: the account is unbound at the end of the call,
// and the caller is at the console now.
func ResetUser(username string) {
	db := database.GetDB()
	var user model.User
	db.Where(&model.User{Username: username}).First(&user)
	user.LastLoginTime = 0
	user.TFASecret = ""
	user.TFABootstrapIP = ""
	user.TFALastUsedStep = 0
	db.Save(&user)

	if err := model.OpenBootstrapWindow(); err != nil {
		logger.Error(err)
	}
}
