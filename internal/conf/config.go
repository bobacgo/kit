package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Config struct {
	mu sync.RWMutex
	c  *config
}

func New(filepath string) *Config {
	cfg := Config{mu: sync.RWMutex{}, c: &config{}}
	cfg.readByFile(filepath)
	version, configFiles := cfg.App().Version, cfg.App().ConfigFileNames
	for _, filename := range configFiles { // 加载多个配置文件
		fullPath := fmt.Sprintf("./deploy/%s/%s", version, filename)
		cfg.readByFile(fullPath)
	}
	cfg.printConfInfo()
	return &cfg
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

func (c *Config) readByFile(path string) {
	_, file, _, ok := runtime.Caller(4)
	if !ok {
		panic("runtime.Caller get path fail")
	}
	abs, _ := filepath.Abs(file)
	fullPath := filepath.Join(filepath.Dir(abs), path)
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	if err = c.Sync(bytes); err != nil {
		panic(err)
	}
}

// 打印读取到的配置信息
func (c *Config) printConfInfo() {
	printConf, _ := yaml.Marshal(c.Config())
	log.Println("======================= load config info ========================")
	fmt.Println(string(printConf))
	log.Println("======================= end config info ========================")
	fmt.Println()
}