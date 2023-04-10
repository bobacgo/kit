package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/_examples/router"
	"github.com/gogoclouds/gogo/app"
	"github.com/gogoclouds/gogo/logger"
)

func main() {
	app.Init()

	r := gin.Default()
	router.LoadRouter(r)
	if err := r.Run(":8080"); err != nil {
		logger.Error(err.Error())
	}
}
