package db

import (
	"fmt"

	"gorm.io/gorm"
)

type InstanceConfig struct {
	Driver gorm.Dialector
	Config Config
}

// 多数据源管理
type DBManager struct {
	sources map[string]*gorm.DB
}

func NewDBManager(cfgKV map[string]InstanceConfig) (*DBManager, error) {
	dbs := make(map[string]*gorm.DB)
	for k, cfg := range cfgKV {
		db, err := NewDB(cfg.Driver, cfg.Config)
		if err != nil {
			return nil, fmt.Errorf("k = %s , init err: %v", k, err)
		}
		dbs[k] = db

	}
	return &DBManager{
		sources: dbs,
	}, nil
}

func (m *DBManager) Default() *gorm.DB {
	return m.sources["default"]
}

func (m *DBManager) Get(k string) *gorm.DB {
	return m.sources[k]
}
