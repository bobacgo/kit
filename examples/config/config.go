package config

import "github.com/bobacgo/kit/app/conf"

func Cfg() Service {
	return conf.GetServiceConf[Service]()
}

type Service struct {
	Admin           Admin          `mapstructure:"admin" yaml:"admin"`
	ErrAttemptLimit int            `mapstructure:"errAttemptLimit" yaml:"errAttemptLimit" mask:""`
	Kafka           conf.Transport `mapstructure:"kafka" yaml:"kafka"`
}

type Admin struct {
	Username string `mapstructure:"username" yaml:"username" mask:""`
	Password string `mapstructure:"password" yaml:"password" mask:""`
}