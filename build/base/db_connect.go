package base

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var (
	DB *gorm.DB
)

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector) (err error) {
	// 使用 gorm.Open 连接数据库
	DB, err = gorm.Open(dbConfig)

	// 处理错误
	if err != nil {
		return err
	}

	if viper.GetString("db.connection") == "mysql" {
		// 获取底层的 SqlDB
		sqlDB, errs := DB.DB()
		if errs != nil {
			return errs
		}
		// 设置最大连接数
		sqlDB.SetMaxOpenConns(viper.GetInt("db.mysql.max_open_connections"))
		// 设置最大空闲连接数
		sqlDB.SetMaxIdleConns(viper.GetInt("db.mysql.max_idle_connections"))
		// 设置每个连接的过期时间
		sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("db.mysql.max_life_seconds")) * time.Second)
	}
	return nil
}
