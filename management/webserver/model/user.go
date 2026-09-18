package model

import (
	"gorm.io/gorm/clause"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

type User struct {
	Base
	Username string `gorm:"uniqueIndex;not null"`
	Password string
	Comment  string

	TFAEnabled    bool   `gorm:"column:tfa_enabled;default:true"`
	TFASecret     string `gorm:"column:tfa_secret"`
	LastLoginTime int64  `gorm:"default:0"`

	// TFALastUsedStep is the time step (unix time divided by the TOTP period)
	// of the passcode that was accepted last. A passcode is only accepted when
	// it belongs to a later step, so a captured passcode cannot be replayed
	// while it is still inside the validity window.
	TFALastUsedStep int64 `gorm:"column:tfa_last_used_step;default:0"`

	IsEnabled bool `gorm:"default:true"`
}

func initAdminUser() error {
	db := database.GetDB()
	user := User{
		Username: constants.SuperUser,
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	return nil
}
