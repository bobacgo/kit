package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/bobacgo/kit/app/types"
	"github.com/coocood/freecache"
)

type freeCache struct {
	cache *freecache.Cache
}

var _ Cache = (*freeCache)(nil)

func NewFreeCache(maxMemorySize types.ByteSize) (Cache, error) {
	if maxMemorySize == "" {
		return DefaultCache(), nil
	}
	size, err := maxMemorySize.ToInt()
	if err != nil {
		return nil, err
	}
	return &freeCache{
		cache: freecache.NewCache(int(size)),
	}, nil
}

func (f *freeCache) SetMaxMemory(_ string) bool {
	return false
}

func (f *freeCache) Set(key string, val any, expire time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewDecoder(&buf).Decode(val); err != nil {
		return err
	}
	return f.cache.Set([]byte(key), buf.Bytes(), int(expire.Seconds()))
}

func (f *freeCache) Get(key string, result any) error {
	value, err := f.cache.Get([]byte(key))
	if err != nil {
		return err
	}
	return gob.NewEncoder(bytes.NewBuffer(value)).Encode(&result)
}

func (f *freeCache) Del(key string) bool {
	ok := f.cache.Del([]byte(key))
	return ok
}

func (f *freeCache) Exists(key string) bool {
	_, err := f.cache.Get([]byte(key))
	return err == nil
}

func (f *freeCache) Clear() bool {
	f.cache.Clear()
	return true
}

func (f *freeCache) Keys() int64 {
	count := f.cache.EntryCount()
	return count
}
