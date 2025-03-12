package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"time"

	"github.com/bobacgo/kit/examples/config"
	_ "github.com/bobacgo/kit/examples/docs"

	kserver "github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/examples/internal/server"

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
		// app.WithMustDB(),
		// app.WithMustRedis(),
		app.WithGinServer(router.Register),
		app.WithGrpcServer(nil),
		app.WithServer(server.KafkaServerName, func(a *app.Options) kserver.Server {
			return new(server.KafkaServer)
		}),
		app.WithAfterStart(func(ctx context.Context, opts *app.Options) error {
			slog.Info("after start")
			return nil
		}),
		app.WithAfterStart(func(ctx context.Context, opts *app.Options) error {
			time.Sleep(3 * time.Second)
			slog.Debug("这是一个debug", "level", "debug")
			slog.Info("这是一个info", "level", "info")
			slog.Warn("这是一个warn", "level", "warn")
			slog.Error("这是一个error", "level", "error")
			return nil
		}),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}
