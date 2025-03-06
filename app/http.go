package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/enum"

	"github.com/bobacgo/kit/app/server/http/middleware"
	"github.com/bobacgo/kit/web/r"
	"github.com/gin-gonic/gin"
)

func RunMustHttpServer(app *App, register func(e *gin.Engine, a *Options)) {
	app.wg.Add(1)
	defer app.wg.Done()

	cfg := app.opts.Conf()

	switch cfg.Env {
	case enum.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	case enum.EnvDev:
		gin.SetMode(gin.DebugMode)
	case enum.EnvTest:
		gin.SetMode(gin.TestMode)
	}

	e := gin.New()
	e.ContextWithFallback = true

	e.Use(gin.Logger())
	e.Use(middleware.Recovery())
	e.Use(middleware.LoggerResponseFail())

	if strings.EqualFold(string(cfg.Env), string(enum.EnvDev)) {
		slog.Warn(fmt.Sprintf(`[gin] Running in "%s" mode`, gin.Mode()))
	}

	healthApi(e, cfg) // provide health API
	pprofApi(e)       // provide pprof API

	if register != nil {
		register(e, &app.opts) // register router
	}

	srv := &http.Server{Addr: cfg.Server.Http.Addr, Handler: e}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		localhost, _ := getRegistryUrl("http", cfg.Server.Http.Addr)
		slog.Info("http server running " + localhost)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("listen: %s\n", err)
		}
	}()

	<-app.exit
	slog.Info("Shutting down http server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Http server forced to shutdown", "err", err)
	}
	slog.Info("http server exiting")
}

// healthApi http check-up API
func healthApi(e *gin.Engine, cfg *conf.Basic) {
	e.GET("/health", func(c *gin.Context) {
		msg := fmt.Sprintf("%s [env=%s] %s, is active", cfg.Name, cfg.Env, cfg.Version)
		r.Reply(c, msg)
	})
}

// profileApi 添加性能分析路由
// TODO 只启动gRPC服务时开启一个http服务 提供性能分析
func pprofApi(e *gin.Engine) {
	// 添加 pprof 路由
	e.GET("/debug/pprof/", gin.WrapF(pprof.Index))
	e.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	e.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	e.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	e.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))
	e.GET("/debug/pprof/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
	e.GET("/debug/pprof/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	e.GET("/debug/pprof/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	e.GET("/debug/pprof/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	e.GET("/debug/pprof/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	e.GET("/debug/pprof/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
}