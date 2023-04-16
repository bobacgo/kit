package main

import (
	"context"
	"flag"
	"github.com/gogoclouds/gogo/_examples/internal/app/module_one/model"

	"github.com/gogoclouds/gogo/_examples/router"
	"github.com/gogoclouds/gogo/app"
)

var config = flag.String("config", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	app.
		New(context.Background(), *config).
		OpenCacheDB().OpenDB(model.Tables).
		CreateHttpServer(router.LoadRouter).
		Run()
}