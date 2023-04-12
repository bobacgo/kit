package main

import (
	"github.com/gogoclouds/gogo/internal/conf"
	"github.com/gogoclouds/gogo/internal/log"
	"github.com/gogoclouds/gogo/logger"
	"testing"
)

func init() {
	log.Init("gogo", conf.Log{
		Level:       "info",
		FileSizeMax: 10,
		FileAgeMax:  10,
	})
}

func TestLogger(t *testing.T) {
	logger.Info("Ok", "test")
}