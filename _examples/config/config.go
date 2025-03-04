package config

var Cfg *Service

type Service struct {
	ErrAttemptLimit int `yaml:"errAttemptLimit"`
}
