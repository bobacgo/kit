package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Configure = new(configure)

type configure struct{}

func (c *configure) ReadFile(path string, conf *Config) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err = conf.Sync(bytes); err != nil {
		panic(err)
	}
}

func (c *configure) PrintConfInfo(conf *Config) {
	// 打印读取到的配置信息
	printConf, _ := yaml.Marshal(conf.Config())
	log.Println("======================= load config info ========================")
	fmt.Println(string(printConf))
}