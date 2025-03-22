package cache

import (
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/maps"
)

const defaultInstanceKey = "default"

// 多数据源管理
type RedisManager map[string]redis.UniversalClient

func NewDBManager(cfgMap map[string]RedisConf) (RedisManager, error) {
	if _, ok := cfgMap[defaultInstanceKey]; !ok {
		return nil, fmt.Errorf("not found default instance, must be has default")
	}
	manager := make(RedisManager, len(cfgMap))
	for k, v := range cfgMap {
		var err error
		if manager[k], err = NewRedis(v); err != nil {
			return nil, fmt.Errorf("k = %s , init err: %v", k, err)
		}
	}

	slog.Info(fmt.Sprintf("[redis] instances object %+q", maps.Keys(cfgMap)))
	return manager, nil
}

func (m RedisManager) Default() redis.UniversalClient {
	return m[defaultInstanceKey]
}

func (m RedisManager) Get(k string) redis.UniversalClient {
	return m[k]
}
