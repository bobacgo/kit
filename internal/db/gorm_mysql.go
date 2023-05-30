package db

import (
	"github.com/gogoclouds/gogo/internal/conf"
	driverMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var mysql = mysqlServer{}

type mysqlServer struct{}

func (mysqlServer) Open(conf conf.Database) (*gorm.DB, error) {
	source := conf.Source
	return gorm.Open(driverMysql.Open(source), &gorm.Config{
		CreateBatchSize:                          1000, // 批量插入每次拆成 1k 条
		QueryFields:                              true, // 会根据当前model的所有字段名称进行 select
		PrepareStmt:                              true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表命名
		},
		DryRun: conf.DryRun,
		Logger: Server.Logger(),
	})
}