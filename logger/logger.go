package logger

import (
	"github.com/gogoclouds/gogo/internal/log"
)

func Info(args ...any) {
	logger.L().Info(args...)
}

func Infof(template string, args ...any) {
	logger.L().Infof(template, args...)
}

func Infow(msg string, keysAndValues ...any) {
	logger.L().Infow(msg, keysAndValues...)
}

func Error(args ...any) {
	logger.L().Error(args...)
}

func Errorf(template string, args ...any) {
	logger.L().Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...any) {
	logger.L().Errorw(msg, keysAndValues...)
}

func Debug(args ...any) {
	logger.L().Debug(args...)
}

func Debugf(template string, args ...any) {
	logger.L().Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...any) {
	logger.L().Debugw(msg, keysAndValues...)
}
