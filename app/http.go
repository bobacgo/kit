package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"

	"errors"

	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/app/validator"
	"github.com/bobacgo/kit/enum"
	pkgvalidator "github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/bobacgo/kit/app/server/http/middleware"
	"github.com/bobacgo/kit/web/r"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type HttpServer struct {
	Opts       *AppOptions
	RegistryFn func(e *gin.Engine, a *AppOptions)

	server *http.Server
}

func NewHttpServer(register func(e *gin.Engine, a *AppOptions), opts *AppOptions) server.Server {
	return &HttpServer{
		Opts:       opts,
		RegistryFn: register,
	}
}

func (srv *HttpServer) Get() any {
	return nil
}

func (srv *HttpServer) Start(ctx context.Context) error {
	cfg := srv.Opts.Conf()

	switch cfg.Env {
	case enum.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	case enum.EnvDev:
		gin.SetMode(gin.DebugMode)
	case enum.EnvTest:
		gin.SetMode(gin.TestMode)
	}

	e := gin.New()
	e.ContextWithFallback = true // 兼容 gin.Context.Get()

	// 扩展 validator 功能
	valid, ok := binding.Validator.Engine().(*pkgvalidator.Validate)
	if ok {
		*valid = *validator.Get()
	}

	e.Use(gin.Logger())
	e.Use(middleware.Recovery())
	e.Use(middleware.LoggerResponseFail())
	if cfg.Otel.Tracer.GrpcEndpoint != "" {
		e.Use(otelgin.Middleware(cfg.Name))
	}

	if strings.EqualFold(string(cfg.Env), string(enum.EnvDev)) {
		slog.Warn(fmt.Sprintf(`[gin] Running in "%s" mode`, gin.Mode()))
	}

	srv.healthApi(e, cfg) // provide health API
	srv.swaggerApi(e)     // provide swagger API
	srv.pprofApi(e)       // provide pprof API

	if srv.RegistryFn != nil {
		srv.RegistryFn(e, srv.Opts) // register router
	}

	// 保证端口监听成功
	listen, err := net.Listen("tcp", cfg.Server.Http.Addr)
	if err != nil {
		return err
	}
	srv.server = &http.Server{Handler: e}

	localhost, _ := getRegistryUrl("http", cfg.Server.Http.Addr)
	slog.Info("[http] http server running " + localhost)
	slog.Info("[http] API docs " + localhost + "/swagger/index.html")
	go func(lit net.Listener) {
		if err := srv.server.Serve(lit); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("listen: %s\n", err)
		}
	}(listen)
	return nil
}

func (srv *HttpServer) Stop(ctx context.Context) error {
	if srv.server == nil {
		return nil
	}
	slog.Info("[http] Shutting down http server...")
	return srv.server.Shutdown(ctx)
}

// healthApi http check-up API
func (srv *HttpServer) healthApi(e *gin.Engine, cfg *conf.Basic) {
	e.GET("/health", func(c *gin.Context) {
		msg := fmt.Sprintf("%s [env=%s] %s, is active", cfg.Name, cfg.Env, cfg.Version)
		r.Reply(c, msg)
	})
}

// swaggerApi swagger 文档
// 访问地址: http://localhost:8080/swagger/index.html
func (srv *HttpServer) swaggerApi(e *gin.Engine) {
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// pprofApi 添加性能分析路由
// TODO 只启动gRPC服务时开启一个http服务 提供性能分析
func (srv *HttpServer) pprofApi(e *gin.Engine) {
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
