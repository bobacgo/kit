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
	"github.com/fsnotify/fsnotify"

	"golang.org/x/exp/maps"

	"github.com/bobacgo/kit/app/registry"
	"github.com/bobacgo/kit/pkg/network"
	"github.com/bobacgo/kit/pkg/uid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	AppOptions

	signal chan os.Signal // 终止服务的信号

	mu       sync.Mutex
	instance *registry.ServiceInstance
}

func New[T any](configPath string, opts ...AppOption) *App {
	// 1. 加载配置
	cfg, err := conf.LoadApp[T](configPath, func(e fsnotify.Event) {
		logger.SetLevel(conf.GetBasicConf().Logger.Level)
		slog.Warn("[config] config onchange", "name", e.Name, "op", e.Op)
	})
	if err != nil {
		log.Panic(err)
	}
	// 2. 初始化日志, 并更新配置
	cfg.Logger = logger.New(cfg.Name, cfg.Logger)

	// 提供一个脱敏标签(mask)的配置文件
	// 扫描标签 并用 ** 替换
	maskConf := tag.Desensitize(cfg)
	cfgData, _ := yaml.Marshal(maskConf)
	slog.Info("[server] local config info\n" + string(cfgData))

	slog.Info(fmt.Sprintf(initDoneFmt, "config"))
	slog.Info(fmt.Sprintf(initDoneFmt, "logger"))

	wg, _ := errgroup.WithContext(context.Background())
	o := AppOptions{
		appId:   uid.UUID(),
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		conf:    conf.GetBasicConf(),
		wg:      wg,
		servers: make(map[string]server.Server),
	}

	// 3. 初始化本地缓存组件
	wg.Go(func() error {
		var err error
		if o.localCache, err = cache.NewFreeCache(o.conf.LocalCache.MaxSize); err != nil {
			return fmt.Errorf("init local cache failed: %w", err)
		}
		components[compCache] = struct{}{}
		slog.Info(fmt.Sprintf(initDoneFmt, compCache))
		return nil
	})
	for _, opt := range opts {
		opt(&o)
	}
	if err := wg.Wait(); err != nil { // 等待 options 实例化结束
		log.Panic(err)
	}

	return &App{
		AppOptions: o,
		signal:     make(chan os.Signal, 1),
	}
}

// Run run server
// 1.注册服务
// 2.退出相关组件或服务
func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for _, fn := range a.beforeStart {
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
	for k, srv := range a.servers {
		a.wg.Go(func() error {
			if err := srv.Start(ctx); err != nil {
				return fmt.Errorf("start %s error: %w", k, err)
			}
			slog.Info(fmt.Sprintf(initDoneFmt, k))
			return nil
		})
	}

	if err = a.wg.Wait(); err != nil { // 等待所有 server 启动完成
		return fmt.Errorf("start servers failed: %w", err)
	}

	if a.registrar != nil {
		timeout := a.Conf().Registry.Timeout
		if timeout != "" {
			ctx, cancel = context.WithTimeout(ctx, timeout.TimeDuration())
			defer cancel()
		}
		if err := a.registrar.Registry(ctx, instance); err != nil {
			return fmt.Errorf("register service failed: %w", err)
		}
	}

	slog.Info("[server] server started")
	slog.Info("[server] app info", "ID", a.AppID(), "name", a.Conf().Name, "version", a.Conf().Version)
	slog.Info(fmt.Sprintf("[server] use components list %q\n", maps.Keys(components)))

	for _, fn := range a.afterStart {
		if err := fn(ctx, &a.AppOptions); err != nil {
			return err
		}
	}

	// 阻塞并监听退出信号
	signal.Notify(a.signal, a.sigs...)
	<-a.signal

	if err = a.shutdown(ctx); err != nil {
		return err
	}
	slog.Info("[server] service has exited")
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
	for _, fn := range a.beforeStop {
		if err := fn(ctx, &a.AppOptions); err != nil {
			return err
		}
	}

	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()

	if a.registrar != nil && instance != nil {
		duration := a.Conf().Registry.Timeout
		if duration == "" {
			duration = "5s"
		}
		ctx, cancel := context.WithTimeout(ctx, duration.TimeDuration())
		defer cancel()
		if err := a.registrar.Deregister(ctx, instance); err != nil {
			return fmt.Errorf("deregister service error: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for k, srv := range a.servers {
		a.wg.Go(func() error {
			if err := srv.Stop(ctx); err != nil { // 出错也不能影响其他服务的停止
				slog.Error("[server] stop server error", k, err)
			} else {
				slog.Info("[server] stop " + k + " service success")
			}
			return nil
		})
	}

	if err := a.wg.Wait(); err != nil {
		return err
	}

	for _, fn := range a.afterStop {
		if err := fn(ctx, &a.AppOptions); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	httpScheme, grpcScheme := false, false
	for _, e := range a.endpoints {
		switch strings.ToLower(e.Scheme) {
		case "https", "http":
			httpScheme = true
		case "grpc":
			grpcScheme = true
		}
		endpoints = append(endpoints, e.String())
	}

	serverCfg := a.Conf().Server

	if !httpScheme && serverCfg.Http.Addr != "" {
		if rUrl, err := getRegistryUrl("http", serverCfg.Http.Addr); err == nil {
			endpoints = append(endpoints, rUrl)
		} else {
			slog.Error("[server] get http registry err", "err", err)
		}
	}
	if !grpcScheme && serverCfg.Rpc.Addr != "" {
		if rUrl, err := getRegistryUrl("grpc", serverCfg.Rpc.Addr); err == nil {
			endpoints = append(endpoints, rUrl)
		} else {
			slog.Error("[server] get grpc registry err", "err", err)
		}
	}
	return &registry.ServiceInstance{
		ID:        a.appId,
		Name:      a.Conf().Name,
		Version:   a.conf.Version,
		Metadata:  nil,
		Endpoints: endpoints,
	}, nil
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
