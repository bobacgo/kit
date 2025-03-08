package conf

import (
	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/db"
	"github.com/bobacgo/kit/app/logger"
	"github.com/bobacgo/kit/app/security"
	"github.com/bobacgo/kit/app/types"
	"github.com/bobacgo/kit/enum"
)

// App 应用配置设计
/*
	1.支持热加载
	2.多种文件类型
	3.配置值格式校验
	4.配置值默认值
	5.配置特殊类型解析
	6.支持配置值脱敏输出
	7.支持多配置文件
	优先级: (相同key)

		1.主配置文件优先级最高
		2.configs 数组索引越小优先级越高
*/
type App struct {
	Basic   `mapstructure:",squash"`
	Service map[string]any `mapstructure:"service"` // 应用自己的其他配置
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
	LocalCache cache.LocalCacheConf `mapstructure:"localCache" yaml:"localCache"`
	Redis      cache.RedisConf      `mapstructure:"redis"`
}

type Transport struct {
	Addr    string         `mapstructure:"addr"`    // 监听地址 0.0.0.0:80
	Timeout types.Duration `mapstructure:"timeout"` // 超时时间 1s
}