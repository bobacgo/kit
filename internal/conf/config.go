package conf

import (
	"gopkg.in/yaml.v3"
	"sync"
)

type Config struct {
	mu sync.RWMutex
	c  *config
}

func New() *Config {
	return &Config{mu: sync.RWMutex{}, c: &config{}}
}

// Sync 局部更新
func (c *Config) Sync(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return yaml.Unmarshal(data, c.c)
}

func (c *Config) App() App {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.c.App
}

func (c *Config) AppServiceKV() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.c.App.ServiceKV
}

func (c *Config) Log() Log {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.c.Log
}

func (c *Config) Database() Database {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.c.Data.Database
}

func (c *Config) Redis() Redis {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.c.Data.Redis
}

func (c *Config) Config() config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return *c.c
}