package cache

import (
	"time"

	"github.com/bobacgo/kit/app/types"
)

// https://zhuanlan.zhihu.com/p/635603181

type Cache interface {
	// SetMaxMemory size : 1KB 100KB 1M 2MB 1GB
	SetMaxMemory(size string) bool

	Set(key string, val any, expire time.Duration) error

	Get(key string, result any) error

	Del(key string) bool

	Exists(key string) bool

	Clear() bool
	// Keys 获取所有缓存中 key 的数量
	Keys() int64
}

var defaultSize types.ByteSize = "512M"

func DefaultCache() (Cache, error) {
	return NewFreeCache(defaultSize)
}
