package db

import "github.com/bobacgo/kit/app/types"

type Config struct {
	Driver        string         `mapstructure:"driver" yaml:"driver"`               // 驱动名称
	Source        string         `mapstructure:"source" mask:":([^@]+)@"`            // root:****@tcp(127.0.0.1:3306)/test
	DryRun        bool           `mapstructure:"dryRun" yaml:"dryRun"`               // 是否为测试模式（空跑sql，不会实际操作数据库）
	SlowThreshold types.Duration `mapstructure:"slowThreshold" yaml:"slowThreshold"` // 慢日志阈值
	MaxOpenConn   int            `mapstructure:"maxOpenConn" yaml:"maxOpenConn"`     // 最大连接数 (高并发 500，低并发 100)
	MaxIdleConn   int            `mapstructure:"maxIdleConn" yaml:"maxIdleConn"`     // 最大空闲连接数 (高并发 50，低并发 10)
	MaxLifeTime   types.Duration `mapstructure:"maxLifeTime" yaml:"maxLifeTime"`     // 最大连接时间 (高并发 1h，低并发 30m)
	MaxIdleTime   types.Duration `mapstructure:"maxIdleTime" yaml:"maxIdleTime"`     // 最大空闲时间 (高并发 15m，低并发 10m)
}
