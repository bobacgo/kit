package logger

import "github.com/gogoclouds/gogo/internal/log"

func Info(args ...interface{}) {
	log.Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	log.Logger.Infof(template, args...)
}

func Error(args ...interface{}) {
	log.Logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	log.Logger.Errorf(template, args...)
}

func Debug(args ...interface{}) {
	log.Logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Logger.Debugf(template, args...)
}