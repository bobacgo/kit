package conf

import (
	"github.com/bobacgo/kit/cache"
	"github.com/bobacgo/kit/db"
	"github.com/bobacgo/kit/enum"
	"github.com/bobacgo/kit/logger"
	"github.com/bobacgo/kit/security"
	"github.com/bobacgo/kit/types"
)

var Cfg Basic

type ServiceConfig[T any] struct {
	Basic   `mapstructure:",squash"`
	Service T `mapstructure:"service"` // 应用自己的其他配置
}

// Basic 服务必要的配置文件
type Basic struct {
	Name    string       `mapstructure:"name" validate:"required"`    // 服务名称
	Version string       `mapstructure:"version" validate:"required"` // 服务版本
	Env     enum.EnvType `mapstructure:"env" validate:"required"`
	// 和主配置文件的在同一个目录可以只写文件名加后缀
	Configs []string `mapstructure:"configs"` // 其他配置文件的路径
	// 注册中心的地址
	Registry Transport `mapstructure:"registry"`
	Server   struct {
		Http Transport `mapstructure:"http"`
		Rpc  Transport `mapstructure:"rpc"` // rpc 端口号没有指定,就是http端口号+1000
	} `mapstructure:"server"`
	Security   security.Config      `mapstructure:"security"`
	Logger     logger.Config        `mapstructure:"logger"`
	DB         map[string]db.Config `mapstructure:"db"` // 支持多数据源 default key 必须存在
	LocalCache cache.LocalCacheConf `mapstructure:"local_cache"`
	Redis      cache.RedisConf      `mapstructure:"redis"`
}

type Transport struct {
	Addr    string         `mapstructure:"addr"`    // 监听地址 0.0.0.0:80
	Timeout types.Duration `mapstructure:"timeout"` // 超时时间 1s
}
