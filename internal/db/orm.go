package db

import (
	"context"
	"fmt"
	"github.com/gogoclouds/gogo/internal/conf"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Server = &server{}

type server struct{}

func (server) NewDB(ctx context.Context, conf *conf.Config) (*gorm.DB, error) {
	driver := conf.Database().Driver
	switch driver {
	case "mysql":
		return mysql.Open(conf.Database())
	default:
		return nil, fmt.Errorf("no dirver: %s", driver)
	}
}

// AutoMigrate create db table
func (server) AutoMigrate(db *gorm.DB, model []any) error {
	return db.AutoMigrate(model...)
}

func (server) Logger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: false,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)
}