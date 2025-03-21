package db

type Config struct {
	Driver        string `mapstructure:"driver" yaml:"driver"`               // 驱动名称
	Source        string `mapstructure:"source" mask:":([^@]+)@"`            // root:****@tcp(127.0.0.1:3306)/test
	DryRun        bool   `mapstructure:"dryRun" yaml:"dryRun"`               // 是否为测试模式（空跑sql，不会实际操作数据库）
	SlowThreshold int    `mapstructure:"slowThreshold" yaml:"slowThreshold"` // 慢日志阈值
	MaxLifeTime   int    `mapstructure:"maxLifeTime" yaml:"maxLifeTime"`     // 最大连接时间
	MaxOpenConn   int    `mapstructure:"maxOpenConn" yaml:"maxOpenConn"`     // 最大连接数
	MaxIdleConn   int    `mapstructure:"maxIdleConn" yaml:"maxIdleConn"`     // 最大空闲连接数
}
