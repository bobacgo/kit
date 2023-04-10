package app

import "github.com/gogoclouds/gogo/logger"

func Init() {
	initLogger()
}

func initLogger() {
	logger.Init("gogo", logger.Config{
		Level:       "debug", // debug | info | error
		FileSizeMax: 10,      // 10 MB
		FileAgeMax:  10,      // 10d
		DirPath:     "/logs",
	})
}
