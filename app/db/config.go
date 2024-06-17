package db

type Config struct {
	Source        string `mapstructure:"source"` // root:root@tcp(127.0.0.1:3306)/test
	DryRun        bool   `mapstructure:"dryRun" yaml:"dryRun"`
	SlowThreshold int    `mapstructure:"slowThreshold" yaml:"slowThreshold"`
	MaxLifeTime   int    `mapstructure:"maxLifeTime" yaml:"maxLifeTime"`
	MaxOpenConn   int    `mapstructure:"maxOpenConn" yaml:"maxOpenConn"`
	MaxIdleConn   int    `mapstructure:"maxIdleConn" yaml:"maxIdleConn"`
}
