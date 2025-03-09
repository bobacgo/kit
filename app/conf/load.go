package conf

import (
	"flag"
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
	// 创建一个新的onChange处理函数,用于处理配置文件的变更
	createOnChange := func(configPath string) func(e fsnotify.Event) {
		return func(e fsnotify.Event) {
			// 配置文件变更时,需要重新加载所有配置
			newCfg := new(App[T])

			// 先加载主配置文件
			if err := Load(filepath, newCfg, nil); err != nil {
				slog.Error("reload main config failed", "error", err)
				return
			}

			// 加载其他配置文件
			// configs 数组索引越小优先级越高
			for i := len(newCfg.Configs) - 1; i >= 0; i-- {
				if err := Load(newCfg.Configs[i], newCfg, nil); err != nil {
					slog.Error("reload sub config failed", "error", err, "config", newCfg.Configs[i])
					return
				}
			}

			// 主配置文件优先级最高,最后加载以覆盖其他配置
			if len(newCfg.Configs) > 0 {
				if err := Load(filepath, newCfg, nil); err != nil {
					slog.Error("reload main config again failed", "error", err)
					return
				}
			}

			if err := validator.Struct(newCfg); err != nil {
				slog.Error("reload config failed", "error", err)
				return
			}

			SetApp(newCfg)
			onChange(e)
		}
	}

	// 加载主配置文件
	if err := Load(filepath, cfg, createOnChange(filepath)); err != nil {
		return nil, err
	}

	// 加载其他配置文件
	// configs 数组索引越小优先级越高
	for i := len(cfg.Configs) - 1; i >= 0; i-- {
		configPath := cfg.Configs[i] // 捕获循环变量
		if err := Load(configPath, cfg, createOnChange(configPath)); err != nil {
			return nil, err
		}
	}

	// 主配置文件优先级最高,最后加载以覆盖其他配置
	if len(cfg.Configs) > 0 {
		if err := Load(filepath, cfg, createOnChange(filepath)); err != nil {
			return nil, err
		}
	}
	if err := validator.Struct(cfg); err != nil {
		return nil, err
	}
	SetApp(cfg)
	return cfg, nil
}

// LoadDefault ./config.yaml
func LoadDefault[T any](onChange func(e fsnotify.Event)) (*T, error) {
	cfgValue := &atomic.Value{}
	cfg := new(T)
	cfgValue.Store(cfg)

	err := Load(".", cfg, func(e fsnotify.Event) {
		newCfg := new(T)
		if err := Load(".", newCfg, onChange); err != nil {
			slog.Error("reload config failed", "error", err)
			return
		}
		cfgValue.Store(newCfg)
		onChange(e)
	})
	return cfgValue.Load().(*T), err
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