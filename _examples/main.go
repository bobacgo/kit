package main

import (
	"flag"
	"github.com/gogoclouds/gogo/_examples/internal/app"
)

var config = flag.String("config", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	app.Run(*config)
}