package base

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConf struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

var db *gorm.DB

func initMysql() {
	var err error
	db, err = gorm.Open(
		mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			Config.MysqlConf.User,
			Config.MysqlConf.Password,
			Config.MysqlConf.Host,
			Config.MysqlConf.Port,
			Config.MysqlConf.Database)),
		&gorm.Config{},
	)
	if err != nil {
		NewLog("fatal", err, "database init failed", "MigrateDB()")
	}
}

func MigrateDB(dbs ...interface{}) {
	if err := db.Set("gorm:table_options", "CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(dbs...); err != nil {
		NewLog("fatal", err, "database migrate failed", "MigrateDB()")
	}
}

func DB() *gorm.DB {
	if Config.MainConfig.Level == logrus.DebugLevel.String() || Config.MainConfig.Level == logrus.TraceLevel.String() {
		return db.Debug()
	}
	return db
}
