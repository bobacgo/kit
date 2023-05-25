package logger

import (
	"github.com/gogoclouds/gogo/internal/log"
)

func Info(args ...interface{}) {
	logger.L().Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.L().Infof(template, args...)
}

func Error(args ...interface{}) {
	logger.L().Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.L().Errorf(template, args...)
}

func Debug(args ...interface{}) {
	logger.L().Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.L().Debugf(template, args...)
}
