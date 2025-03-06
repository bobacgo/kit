package config

var Cfg = new(Service)

type Service struct {
	ErrAttemptLimit int `yaml:"errAttemptLimit"`
}