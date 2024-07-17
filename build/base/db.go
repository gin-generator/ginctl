package base

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() {
	var dbConfig gorm.Dialector
	switch viper.GetString("db.connection") {
	case "mysql":
		// 构建 DSN 信息
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
			viper.GetString("db.mysql.username"),
			viper.GetString("db.mysql.password"),
			viper.GetString("db.mysql.host"),
			viper.GetInt("db.mysql.port"),
			viper.GetString("db.mysql.database"),
			viper.GetString("db.mysql.charset"),
		)
		dbConfig = mysql.New(mysql.Config{
			DSN:                       dsn,
			SkipInitializeWithVersion: viper.GetBool("db.mysql.skip_initialize_with_version"),
		})
	case "tidb":
		// 构建 DSN 信息
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&tls=%s",
			viper.GetString("db.tidb.username"),
			viper.GetString("db.tidb.password"),
			viper.GetString("db.tidb.host"),
			viper.GetString("db.tidb.port"),
			viper.GetString("db.tidb.database"),
			viper.GetBool("db.tidb.ssl"),
		)
		dbConfig = mysql.New(mysql.Config{
			DSN: dsn,
		})
	//case "sqlite":
	//	dbConfig = sqlite.Open(config.Get("db.sqlite.database"))
	default:
		panic(errors.New("database connection not supported"))
	}

	// 连接数据库，并设置 GORM 的日志模式
	err := Connect(dbConfig)
	if err != nil {
		panic(errors.New("database connection failure"))
	}
}
