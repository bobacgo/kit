package conf

import (
	"flag"
	"github.com/bobacgo/kit/pkg/tag"
	"log/slog"
	"sync/atomic"

	"github.com/bobacgo/kit/app/validator"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	basicCfg   atomic.Value
	serviceCfg atomic.Value
)

func GetBasicConf() Basic {
	cfg, _ := basicCfg.Load().(Basic)
	return cfg
}

func GetServiceConf[T any]() T {
	v, _ := serviceCfg.Load().(T)
	return v
}

func SetApp[T any](appCfg *App[T]) {
	if appCfg == nil {
		return
	}
	basicCfg.Store(appCfg.Basic)
	serviceCfg.Store(appCfg.Service)
}

// LoadApp 加载配置文件
// 配置文件有变化时,会自动全部重新加载配置文件
// 优先级: (相同key)
//
//	1.主配置文件优先级最高
//	2.configs 数组索引越小优先级越高
func LoadApp[T any](filepath string, onChange func(e fsnotify.Event)) (*App[T], error) {
	cfg := new(App[T])

	if onChange != nil {
		onChange = reload[T](filepath, onChange)
	}

	// 加载主配置文件
	if err := Load(filepath, cfg, onChange); err != nil {
		return nil, err
	}

	// 加载其他配置文件
	// configs 数组索引越小优先级越高
	for i := len(cfg.Configs) - 1; i >= 0; i-- {
		configPath := cfg.Configs[i] // 捕获循环变量
		if err := Load(configPath, cfg, onChange); err != nil {
			return nil, err
		}
	}

	// 主配置文件优先级最高,最后加载以覆盖其他配置
	if len(cfg.Configs) > 0 {
		if err := Load(filepath, cfg, nil); err != nil {
			return nil, err
		}
	}
	if err := validator.Struct(cfg); err != nil {
		return nil, err
	}
	cfg = tag.Default(cfg) // 带有默认值 tag 标签赋值
	SetApp(cfg)
	return cfg, nil
}

func reload[T any](path string, onChange func(e fsnotify.Event)) func(e fsnotify.Event) {
	return func(e fsnotify.Event) {
		if _, err := LoadApp[T](path, nil); err != nil {
			slog.Error("reload config error", "err", err)
			return
		}
		if onChange != nil {
			onChange(e)
		}
	}
}

func Load[T any](filepath string, cfg *T, onChange func(e fsnotify.Event)) error {
	vpr := viper.New()
	vpr.SetConfigFile(filepath)
	vpr.ReadInConfig()
	if err := vpr.ReadInConfig(); err != nil {
		return err
	}
	if err := vpr.Unmarshal(cfg); err != nil {
		return err
	}
	if onChange != nil {
		vpr.WatchConfig()
		vpr.OnConfigChange(func(e fsnotify.Event) {
			onChange(e)
		})
	}
	return nil
}

func BindPFlags() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}