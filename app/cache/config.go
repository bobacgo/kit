package cache

import (
	"github.com/bobacgo/kit/app/types"
)

type RedisConf struct {
	Addrs        []string       `mapstructure:"addrs"` // [127.0.0.1:6379, 127.0.0.1:7000]
	Username     string         `mapstructure:"username"`
	Password     string         `mapstructure:"password" mask:""`
	DB           uint8          `mapstructure:"db"`
	PoolSize     int            `mapstructure:"poolSize" yaml:"poolSize"`
	ReadTimeout  types.Duration `mapstructure:"readTimeout" yaml:"readTimeout" validate:"duration"`   // 0.2s
	WriteTimeout types.Duration `mapstructure:"writeTimeout" yaml:"writeTimeout" validate:"duration"` // 0.2s
}

type LocalCacheConf struct {
	MaxSize types.ByteSize `mapstructure:"maxSize" yaml:"maxSize"` // 最大容量
}