package main

import (
	"github.com/gogoclouds/gogo/logger"
	"testing"
)

func init() {
	logger.Init("gogo", logger.Config{
		Level:       "info",
		FileSizeMax: 10,
		FileAgeMax:  10,
	})
}

func TestLogger(t *testing.T) {
	logger.Info("Ok", "test")
}