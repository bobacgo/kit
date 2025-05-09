package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/bobacgo/kit/app/mq/kafka"
	"github.com/bobacgo/kit/app/otel"
	"github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/app/server/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"golang.org/x/sync/errgroup"

	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/db"
	"github.com/bobacgo/kit/app/registry"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	compCache   = "cache"
	compRedis   = "redis"
	compHttp    = "http"
	compRpc     = "rpc"
	compKafka   = "kafka"
	compGateway = "gateway"
)

const initDoneFmt = " [%s] init done."

var components = make(map[string]struct{}) // 组件名列表

type AppOption func(o *AppOptions)

type AppOptions struct {
	appId string      // 应用程序启动实例ID
	sigs  []os.Signal // 监听的程序退出信号

	conf conf.Basic

	wg *errgroup.Group // 用于并发初始化组件
	// 内置功能
	localCache cache.Cache
	redis      cache.RedisManager
	db         db.DBManager

	// hook func
	beforeStart                       []func(ctx context.Context) error
	afterStart, beforeStop, afterStop []func(ctx context.Context, opts *AppOptions) error

	endpoints []*url.URL
	registrar registry.ServiceRegistrar

	// 插件功能 如 服务需要依赖 Kafka 等常驻进程服务
	servers map[string]server.Server
}

// AppID 获取应用程序启动实例ID
func (o *AppOptions) AppID() string {
	return o.appId
}

// Conf 获取公共配置(eg app info、logger config、db config 、redis config)
func (o *AppOptions) Conf() *conf.Basic {
	basicConf := conf.GetBasicConf()
	return &basicConf
}

// LocalCache 获取本地缓存
// 1.一级缓存 变动小、容量少。容量固定，有淘汰策略。
// 2.不适合分布式数据共享。
func (o *AppOptions) LocalCache() cache.Cache {
	return o.localCache
}

// DB 获取数据库连接
// DB gorm 关系型数据库 -- 持久化
func (o *AppOptions) DB() db.DBManager {
	return o.db
}

// Redis 获取redis client
// CacheDB 二级缓存 容量大，有网络IO延迟
func (o *AppOptions) Redis() cache.RedisManager {
	return o.redis
}

// Server 获取自己注入的组件服务
func (o *AppOptions) Server(name string) (any, bool) {
	srv, ok := o.servers[name]
	if !ok {
		return nil, false
	}
	return srv.Get(), true
}

// WithAppID 设置appID，默认UUID
func WithAppID(id string) AppOption {
	return func(o *AppOptions) {
		if id != "" {
			o.appId = id
		}
	}
}

// WithSignal 设置可监听的退出信号
func WithSignal(sigs []os.Signal) AppOption {
	return func(o *AppOptions) {
		if len(sigs) > 0 {
			o.sigs = sigs
		}
	}
}

// WithEndpoints 设置服务地址
func WithEndpoints(endpoints []*url.URL) AppOption {
	return func(o *AppOptions) {
		if len(endpoints) > 0 {
			o.endpoints = endpoints
		}
	}
}

// WithRegistrar 设置服务注册器
func WithRegistrar(registrar registry.ServiceRegistrar) AppOption {
	return func(o *AppOptions) {
		o.registrar = registrar
	}
}

// WithMustRedis 初始化 Redis 组件（错误直接panic）
func WithMustRedis() AppOption {
	components[compRedis] = struct{}{}
	return func(o *AppOptions) {
		o.wg.Go(func() error {
			var err error
			if o.redis, err = cache.NewRedisManager(o.conf.Name, o.conf.Redis); err != nil {
				return fmt.Errorf("init redis failed: %w", err)
			}
			slog.Info(fmt.Sprintf(initDoneFmt, compRedis))
			return nil
		})
	}
}

// WithMustDB 初始化数据库组件（错误直接panic）
// GORM 官方支持的数据库类型有
// MySQL, PostgreSQL, SQLite, SQL Server 和 TiDB
func WithMustDB(drivers ...db.DriverOpenFunc) AppOption {
	components[db.ComponentName] = struct{}{}
	return func(o *AppOptions) {
		o.wg.Go(func() error {
			var err error
			dmap := db.DialectorMap(drivers, o.conf.DB)
			if o.db, err = db.NewDBManager(dmap); err != nil {
				return fmt.Errorf("init db manager failed: %w", err)
			}
			slog.Info(fmt.Sprintf(initDoneFmt, db.ComponentName))
			return nil
		})
	}
}

// WithServer 注入 server， 需要指定唯一的key
func WithServer(name string, srv func(a *AppOptions) server.Server) AppOption {
	components[name] = struct{}{}
	return func(o *AppOptions) {
		if srv != nil {
			o.servers[name] = srv(o)
		}
	}
}

// WithGinServer 使用 gin server
// srv 注册路由，添加 handler
func WithGinServer(svr func(e *gin.Engine, a *AppOptions)) AppOption {
	return WithServer(compHttp, func(a *AppOptions) server.Server {
		return NewHttpServer(svr, a)
	})
}

// WithGrpcServer 使用 grpc server
// srv 注册业务服务（pb生成）
// grpcServerOpts 可选的 grpc server options
func WithGrpcServer(svr func(s *grpc.Server, a *AppOptions), grpcServerOpts ...grpc.ServerOption) AppOption {
	return WithServer(compRpc, func(a *AppOptions) server.Server {
		return NewRpcServer(svr, a, grpcServerOpts...)
	})
}

// WithKafka 使用 Kafka server
// subs 注册消息处理器
func WithKafka(subs ...kafka.Subscriber) AppOption {
	return WithServer(compKafka, func(a *AppOptions) server.Server {
		return kafka.New(a.Conf().Name, a.Conf().Kafka, subs...)
	})
}

// WithCache 使用GrpcGateway
// reg 注册服务
// grpcServerOpts 可选的 grpc server options
func WithGatewayServer(handlerFunc func(mux *runtime.ServeMux), muxs ...runtime.ServeMuxOption) AppOption {
	return WithServer(compGateway, func(a *AppOptions) server.Server {
		return gateway.New(a.Conf().GrpcGateway, handlerFunc, muxs...)
	})
}

// WithTracerServer 使用 TracerServer
func WithTracerServer() AppOption {
	return WithServer("tracer", func(a *AppOptions) server.Server {
		appInfo := otel.AppInfo{
			Name:    a.Conf().Name,
			ID:      a.appId,
			Version: a.Conf().Version,
		}
		return otel.NewTracerServer(appInfo, &a.Conf().Otel.Tracer)
	})
}

// WithBeforeStart 在启动前执行（可以传多次）
func WithBeforeStart(fn func(ctx context.Context) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.beforeStart = append(o.beforeStart, fn)
		}
	}
}

// WithAfterStart 在启动后执行（可以传多次）
func WithAfterStart(fn func(ctx context.Context, opts *AppOptions) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.afterStart = append(o.afterStart, fn)
		}
	}
}

// WithBeforeStop 在停止前执行（可以传多次）
func WithBeforeStop(fn func(ctx context.Context, opts *AppOptions) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.beforeStop = append(o.beforeStop, fn)
		}
	}
}

// WithAfterStop 在停止后执行（可以传多次）
func WithAfterStop(fn func(ctx context.Context, opts *AppOptions) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.afterStop = append(o.afterStop, fn)
		}
	}
}
