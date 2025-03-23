package db

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// NewDB 新建数据库连接
// 内部已自动 ping
func NewDB(dialector gorm.Dialector, conf Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{
		CreateBatchSize:                          1000, // 批量插入每次拆成 1k 条
		QueryFields:                              true, // 会根据当前model的所有字段名称进行 select
		PrepareStmt:                              true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		SkipDefaultTransaction:                   true, // 禁用默认事务操作，提高运行速度
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		TranslateError:                           true, // 启用翻译错误
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表命名
		},
		DryRun: conf.DryRun,
		Logger: Logger(conf),
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open db err: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get DB err: %w", err)
	}

	// 影响最大并发数。
	// 过大可能导致数据库负载过高，过小会限制并发性能。
	//一般设置在 100~500，具体根据数据库负载情况调整。
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn) // 设置最大连接数
	// 控制保持在池中的空闲连接数。
	// 过大会浪费资源，过小可能导致频繁创建连接，增加延迟。
	// 典型范围是 10~50。
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn) // 设置闲置连接数
	// 控制连接存活的最大时间，避免连接长时间占用资源导致 MySQL 关闭连接。
	// 建议设置 30min ~ 1h，防止连接泄露。
	sqlDB.SetConnMaxLifetime(conf.MaxLifeTime.TimeDuration()) // 连接的最大可复用时间
	// 控制空闲连接的最长时间，防止长期空闲的连接占用资源。
	// 典型值 10min，根据业务需求调整。
	sqlDB.SetConnMaxIdleTime(conf.MaxIdleTime.TimeDuration()) // 空闲连接的最大生存时间
	return db, nil
}

func Logger(conf Config) logger.Interface {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks
	if conf.SlowThreshold != "" {
		logger.SlowThreshold = conf.SlowThreshold.TimeDuration() // 慢 SQL 阈值
	}
	return logger
}
