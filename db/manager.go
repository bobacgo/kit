package db

import (
	"log"

	"gorm.io/gorm"
)

type InstanceConfig struct {
	Driver gorm.Dialector
	Config Config
}

// 多数据源管理
type SourceManager struct {
	sources map[string]*gorm.DB
}

func NewSourceManager(cfgKV map[string]InstanceConfig) *SourceManager {
	dbs := make(map[string]*gorm.DB)
	for k, cfg := range cfgKV {
		db, err := NewDB(cfg.Driver, cfg.Config)
		if err != nil {
			log.Panicf("k = %s , init err: %v\n", k, err)
		}
		dbs[k] = db

	}
	return &SourceManager{
		sources: dbs,
	}
}

func (m *SourceManager) Default() *gorm.DB {
	return m.sources["default"]
}

func (m *SourceManager) Get(k string) *gorm.DB {
	return m.sources[k]
}
