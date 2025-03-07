package main

import (
	"flag"
	kserver "github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/examples/internal/server"
	"log"

	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/examples/config"
	"github.com/bobacgo/kit/examples/internal/app/router"
)

var filepath = flag.String("config", "./config.yaml", "config file path")

func init() {
	flag.String("name", "admin-service", "service name")
	flag.String("env", "dev", "run config context")
	flag.String("logger.level", "info", "logger level")
	flag.Int("port", 8080, "http port 8080, rpc port 9080")
	conf.BindPFlags()
}

func main() {
	newApp := app.New(*filepath,
		app.WithScanConfig(config.Cfg),
		app.WithLogger(),
		app.WithLocalCache(),
		// app.WithMustDB(),
		// app.WithMustRedis(),
		app.WithGinServer(router.Register),
		app.WithServer(server.KafkaServerName, func(a *app.Options) kserver.Server {
			return new(server.KafkaServer)
		}),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}