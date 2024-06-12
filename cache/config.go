package cache

import (
	"github.com/bobacgo/kit/types"
)

type RedisConf struct {
	Addrs        []string       `mapstructure:"addrs"` // [127.0.0.1:6379, 127.0.0.1:7000]
	Username     string         `mapstructure:"username"`
	Password     string         `mapstructure:"password"`
	DB           uint8          `mapstructure:"db"`
	PoolSize     int            `mapstructure:"poolSize"`
	ReadTimeout  types.Duration `mapstructure:"readTimeout"`  // 0.2s
	WriteTimeout types.Duration `mapstructure:"writeTimeout"` // 0.2s
}

type LocalCacheConf struct {
	MaxSize types.ByteSize `mapstructure:"maxSize"` // 最大容量
}
