package model

import (
	"CryptoCurrency/common"
	"CryptoCurrency/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

var mysqlDB *gorm.DB

func init() {
	config := config.Configuration.Mysql

	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别为 Silent，来关闭日志
	})
	if err != nil {
		slog.Error("mysql init error:", err.Error())
		common.SystemClose()
	}
	mysqlDB = db
}

func GetMysqlDB() *gorm.DB {
	return mysqlDB
}
