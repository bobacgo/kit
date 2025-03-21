package db

import (
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const ComponentName = "database"

func withPrefix(prefix, format string, msgs ...any) string {
	format = "[" + prefix + "] " + format
	return fmt.Sprintf(format, msgs...)
}

const defaultInstanceKey = "default"

// 多数据源管理
type DBManager struct {
	sources map[string]*gorm.DB
}

func NewDBManager(cfgKV map[string]DialectorConfig) (*DBManager, error) {
	dbs := make(map[string]*gorm.DB)
	for k, cfg := range cfgKV {
		db, err := NewDB(cfg.Dialector, cfg.Config)
		if err != nil {
			return nil, fmt.Errorf("k = %s , init err: %v", k, err)
		}
		dbs[k] = db

	}

	slog.Info(withPrefix(ComponentName, "instances object %+q", maps.Keys(dbs)))

	return &DBManager{
		sources: dbs,
	}, nil
}

func (m *DBManager) Default() *gorm.DB {
	return m.sources[defaultInstanceKey]
}

func (m *DBManager) Get(k string) *gorm.DB {
	return m.sources[k]
}

type DialectorConfig struct {
	Dialector gorm.Dialector
	Config    Config
}

type DriverOpenFunc func(dsn string) gorm.Dialector

func DialectorMap(drivers []DriverOpenFunc, cfgMap map[string]Config) map[string]DialectorConfig {
	dialectorMap := make(map[string]DialectorConfig, len(drivers))

	driverMap := map[string]DriverOpenFunc{mysql.DefaultDriverName: mysql.Open} // 默认提供mysql驱动
	for _, d := range drivers {
		name := d("").Name()
		driverMap[name] = d
	}

	slog.Info(withPrefix(ComponentName, "support driver %+q", maps.Keys(driverMap)))

	for k, c := range cfgMap {
		openFunc, ok := driverMap[c.Driver]
		if !ok {
			slog.Warn(withPrefix(ComponentName, "driver not found, Please check the configuration file"), "driver", c.Driver)
			continue
		}
		dialectorMap[k] = DialectorConfig{Dialector: openFunc(c.Source), Config: c}
	}
	return dialectorMap
}
