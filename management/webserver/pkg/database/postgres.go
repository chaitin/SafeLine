package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
)

type PostgresDB struct {
	*gorm.DB
}

var (
	// db is not designed to be used directly, use method `GetDB` instead.
	db PostgresDB
)

func GetDB() *PostgresDB {
	return &db
}

func InitDB() error {
	dbConfig := config.GlobalConfig.DB

	if dbConfig.URL == "" {
		return errors.New("empty database url")
	}
	// URL also works, see https://godoc.org/github.com/lib/pq
	var (
		gormDB *gorm.DB
		err    error
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "mgt_", // table name prefix, table for `User` would be `t_users`
			SingularTable: true,   // use singular table name, table for `User` would be `user` with this option enabled
		},
		AllowGlobalUpdate: false,
		Logger:            logger.Default.LogMode(logger.Silent),
	}
	if dbConfig.LogSQL {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	gormDB, err = gorm.Open(postgres.Open(dbConfig.URL), gormConfig)

	if err != nil {
		return err
	}

	db.DB = gormDB

	return nil
}

func (db *PostgresDB) SetDB(gormDB *gorm.DB) {
	db.DB = gormDB
}
