package cmd

import (
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

// ResetUser clears the credentials of a user so that the next visit of the
// console binds a new TFA secret.
//
// The pending secret and the address it was issued to are cleared together:
// a secret that was handed to somebody else cannot be bound by the
// administrator, and starting the bootstrap over from the shell is the only
// way out of that state.
func ResetUser(username string) {
	db := database.GetDB()
	var user model.User
	db.Where(&model.User{Username: username}).First(&user)
	user.LastLoginTime = 0
	user.TFASecret = ""
	user.TFABootstrapIP = ""
	user.TFALastUsedStep = 0
	db.Save(&user)
}
