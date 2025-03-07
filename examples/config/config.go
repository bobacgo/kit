package config

import "github.com/bobacgo/kit/app/conf"

var Cfg = new(Service)

type Service struct {
	ErrAttemptLimit int            `mapstructure:"errAttemptLimit" yaml:"errAttemptLimit"`
	Kafka           conf.Transport `mapstructure:"kafka" yaml:"kafka"`
}