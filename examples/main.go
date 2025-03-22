package main

import (
	"context"
	"flag"
	"log"
	"log/slog"

	kserver "github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/examples/internal/config"
	"github.com/bobacgo/kit/examples/internal/server"
	"gorm.io/driver/sqlite"

	_ "github.com/bobacgo/kit/examples/docs"

	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/examples/internal/app/router"
)

var filepath = flag.String("config", "./examples/config.yaml", "config file path")

func init() {
	flag.String("name", "admin-service", "service name")
	flag.String("env", "dev", "run config context")
	flag.String("logger.level", "info", "logger level")
	flag.Int("port", 8080, "http port 8080, rpc port 9080")
	conf.BindPFlags()
}

//go:generate swag init --parseDependency --parseInternal --dir ./ --output ./docs
func main() {
	newApp := app.New[config.Service](*filepath,
		app.WithMustDB(sqlite.Open),
		app.WithMustRedis(),
		// app.WithKafka(),
		app.WithGinServer(router.Register),
		app.WithGrpcServer(server.GrpcRegisterServer),
		app.WithServer(server.JobServerName, func(a *app.AppOptions) kserver.Server {
			return new(server.JobServer)
		}),
		app.WithAfterStart(func(ctx context.Context, opts *app.AppOptions) error {
			slog.Info("after start")
			return nil
		}),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}
