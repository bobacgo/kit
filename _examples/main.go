package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/_examples/router"
	"github.com/gogoclouds/gogo/app"
	"github.com/gogoclouds/gogo/logger"
)

var config = flag.String("config", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	app.Init(*config)

	r := gin.Default()
	router.LoadRouter(r)
	if err := r.Run(":8080"); err != nil {
		logger.Error(err.Error())
	}
}
