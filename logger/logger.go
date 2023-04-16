package logger

import (
	"github.com/gogoclouds/gogo/internal/log"
)

func Info(args ...interface{}) {
	logger.Log.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Log.Infof(template, args...)
}

func Error(args ...interface{}) {
	logger.Log.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Log.Errorf(template, args...)
}

func Debug(args ...interface{}) {
	logger.Log.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Log.Debugf(template, args...)
}