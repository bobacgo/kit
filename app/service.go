package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bobacgo/kit/app/logger"
	"github.com/bobacgo/kit/pkg/tag"
	"gopkg.in/yaml.v2"

	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/app/types"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/bobacgo/kit/app/registry"
	"github.com/bobacgo/kit/pkg/network"
	"github.com/bobacgo/kit/pkg/uid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts Options

	wg     *errgroup.Group
	signal chan os.Signal // 终止服务的信号

	mu       sync.Mutex
	instance *registry.ServiceInstance
}

func New[T any](configPath string, opts ...Option) *App {
	// 加载配置
	cfg, err := conf.LoadApp[T](configPath, func(e fsnotify.Event) {
		logger.SetLevel(conf.GetBasicConf().Logger.Level)
		slog.Warn("config onchange", "name", e.Name, "op", e.Op)
	})
	if err != nil {
		log.Panic(err)
	}
	// 初始化日志配置
	cfg.Logger = newLogger(cfg.Name, cfg.Logger)

	// 提供一个脱敏标签(mask)的配置文件
	// 扫描标签 并用 ** 替换
	maskConf := tag.Desensitize(cfg)
	cfgData, _ := yaml.Marshal(maskConf)
	slog.Info("local config info\n" + string(cfgData))

	slog.Info(fmt.Sprintf(initDoneFmt, "config"))
	slog.Info(fmt.Sprintf(initDoneFmt, "logger"))

	wg, _ := errgroup.WithContext(context.Background())
	o := Options{
		appId:   uid.UUID(),
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		conf:    conf.GetBasicConf(),
		wgInit:  wg,
		servers: make(map[string]server.Server),
	}
	wg.Go(func() error {
		var err error
		o.localCache, err = newLocalCache(o.conf.LocalCache.MaxSize)
		return err
	})
	for _, opt := range opts {
		opt(&o)
	}
	if err := wg.Wait(); err != nil { // 等待 options 实例化结束
		log.Panic(err)
	}

	return &App{
		opts:   o,
		wg:     wg,
		signal: make(chan os.Signal, 1),
	}
}

// Run run server
// 1.注册服务
// 2.退出相关组件或服务
func (a *App) Run() error {
	opts := a.opts

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for _, fn := range opts.beforeStart {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	// start servers
	for k, srv := range opts.servers {
		a.wg.Go(func() error {
			if err := srv.Start(ctx); err != nil {
				return errors.Wrapf(err, "start %s error", k)
			}
			slog.Info(fmt.Sprintf(initDoneFmt, k))
			return nil
		})
	}

	if err = a.wg.Wait(); err != nil { // 等待所有 server 启动完成
		return errors.Wrap(err, "start servers failed")
	}

	if opts.registrar != nil {
		timeout := opts.conf.Registry.Timeout
		if timeout != "" {
			ctx, cancel = context.WithTimeout(ctx, timeout.TimeDuration())
			defer cancel()
		}
		if err := opts.registrar.Registry(ctx, instance); err != nil {
			return errors.Wrap(err, "register service failed")
		}
	}

	slog.Info("server started")
	slog.Info("app info", "ID", opts.AppID(), "name", opts.Conf().Name, "version", opts.Conf().Version)
	slog.Info(fmt.Sprintf("use components list %q\n", maps.Keys(components)))

	for _, fn := range opts.afterStart {
		if err := fn(ctx, &opts); err != nil {
			return err
		}
	}

	// 阻塞并监听退出信号
	signal.Notify(a.signal, opts.sigs...)
	<-a.signal

	if err = a.shutdown(ctx); err != nil {
		return err
	}
	slog.Info("service has exited")
	return nil
}

// Stop 手动停止服务
func (a *App) Stop() {
	a.signal <- os.Interrupt
}

// shutdown server
// 1.注销服务
// 2.退出 http、grpc服务
func (a *App) shutdown(ctx context.Context) error {
	opts := a.opts

	for _, fn := range opts.beforeStop {
		if err := fn(ctx, &opts); err != nil {
			return err
		}
	}

	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()

	if opts.registrar != nil && instance != nil {
		duration := opts.conf.Registry.Timeout
		if duration == "" {
			duration = "5s"
		}
		ctx, cancel := context.WithTimeout(ctx, duration.TimeDuration())
		defer cancel()
		if err := opts.registrar.Deregister(ctx, instance); err != nil {
			return errors.Wrap(err, "deregister service error")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for k, srv := range opts.servers {
		a.wg.Go(func() error {
			if err := srv.Stop(ctx); err != nil { // 出错也不能影响其他服务的停止
				slog.Error("stop server error", k, err)
			} else {
				slog.Info("stop " + k + " service success")
			}
			return nil
		})
	}

	if err := a.wg.Wait(); err != nil {
		return err
	}

	for _, fn := range opts.afterStop {
		if err := fn(ctx, &opts); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	opts := a.opts

	endpoints := make([]string, 0)
	httpScheme, grpcScheme := false, false
	for _, e := range opts.endpoints {
		switch strings.ToLower(e.Scheme) {
		case "https", "http":
			httpScheme = true
		case "grpc":
			grpcScheme = true
		}
		endpoints = append(endpoints, e.String())
	}

	if !httpScheme && opts.conf.Server.Http.Addr != "" {
		if rUrl, err := getRegistryUrl("http", opts.conf.Server.Http.Addr); err == nil {
			endpoints = append(endpoints, rUrl)
		} else {
			slog.Error("get http registry err", "err", err)
		}
	}
	if !grpcScheme && opts.conf.Server.Rpc.Addr != "" {
		if rUrl, err := getRegistryUrl("grpc", opts.conf.Server.Rpc.Addr); err == nil {
			endpoints = append(endpoints, rUrl)
		} else {
			slog.Error("get grpc registry err", "err", err)
		}
	}
	return &registry.ServiceInstance{
		ID:        opts.appId,
		Name:      opts.conf.Name,
		Version:   opts.conf.Version,
		Metadata:  nil,
		Endpoints: endpoints,
	}, nil
}

func newLogger(appName string, logCfg logger.Config) logger.Config {
	if logCfg.Level == "" {
		logCfg.Level = logger.LogLevel_Info
	}
	opts := []logger.Option{
		logger.WithLevel(logCfg.Level),
	}
	if logCfg.TimeFormat != "" {
		opts = append(opts, logger.WithTimeFormat(logCfg.TimeFormat))
	}
	if appName != "" {
		opts = append(opts, logger.WithFilename(appName))
	}
	if logCfg.Filepath != "" {
		opts = append(opts, logger.WithFilepath(logCfg.Filepath))
	}
	if logCfg.FilenameSuffix != "" {
		opts = append(opts, logger.WithFilenameSuffix(logCfg.FilenameSuffix))
	}
	if logCfg.FileExtension != "" {
		opts = append(opts, logger.WithFileExtension(logCfg.FileExtension))
	}
	if logCfg.FileMaxSize > 0 {
		opts = append(opts, logger.WithFileMaxSize(logCfg.FileMaxSize))
	}
	if logCfg.FileMaxAge > 0 {
		opts = append(opts, logger.WithFileMaxAge(logCfg.FileMaxAge))
	}
	if logCfg.FileJsonEncoder {
		opts = append(opts, logger.WithFileJsonEncoder(logCfg.FileJsonEncoder))
	}
	if logCfg.FileCompress {
		opts = append(opts, logger.WithFileCompress(logCfg.FileCompress))
	}

	cfg := logger.NewConfig(opts...)
	// 初始化日志配置
	logger.InitZapLogger(cfg)
	return cfg
}

func newLocalCache(limit types.ByteSize) (cache.Cache, error) {
	if limit == "" {
		return cache.DefaultCache(), nil // 没有指定大小就使用默认的 512MB
	}
	localCache, err := cache.NewFreeCache(limit)
	if err != nil {
		return nil, errors.Wrap(err, "init local cache failed")
	}
	slog.Info(fmt.Sprintf(initDoneFmt, compCache))

	components[compCache] = struct{}{}
	return localCache, nil
}

func getRegistryUrl(scheme, addr string) (string, error) {
	ip, err := network.OutBoundIP()
	if err != nil {
		return "", err
	}
	_, ports, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return scheme + "://" + net.JoinHostPort(ip, ports), nil
}
