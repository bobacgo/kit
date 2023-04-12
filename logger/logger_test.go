package logger_test

import (
	"github.com/gogoclouds/gogo/internal/conf"
	"github.com/gogoclouds/gogo/internal/log"
	"github.com/gogoclouds/gogo/logger"
	"testing"
)

func init() {
	log.Init("gogo", conf.Log{
		Level:       "debug", // debug | info | error
		FileSizeMax: 10,      // 10 MB
		FileAgeMax:  10,      // 10d
		DirPath:     "/logs",
	})
}

func TestLogger(t *testing.T) {
	logger.Debug("The is ", "Debug")
	logger.Info("The is ", "Info")
	logger.Error("The is ", "Error")

	logger.Debugf("The is %s", "Debugf")
	logger.Info("The is %s", "Info")
	logger.Errorf("The is %s", "Errorf")
}