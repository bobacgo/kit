package logger_test

import (
	"testing"

	"github.com/gogoclouds/gogo/internal/conf"
	logging "github.com/gogoclouds/gogo/internal/log"
	"github.com/gogoclouds/gogo/logger"
)

func init() {
	logging.Initialize("gogo", conf.Log{
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
	logger.Infof("The is %s", "Info")
	logger.Errorf("The is %s", "Errorf")
}
