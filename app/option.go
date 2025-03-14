package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/bobacgo/kit/app/server"
	"github.com/pkg/errors"

	"golang.org/x/sync/errgroup"

	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/db"
	"github.com/bobacgo/kit/app/registry"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
)

const (
	compCache = "local_cache"
	compRedis = "redis"
	compDB    = "db"
	compHttp  = "http"
	compRpc   = "rpc"
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
	redis      redis.UniversalClient
	db         *db.DBManager

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
func (o *AppOptions) DB() *db.DBManager {
	return o.db
}

// Redis 获取redis client
// CacheDB 二级缓存 容量大，有网络IO延迟
func (o *AppOptions) Redis() redis.UniversalClient {
	return o.redis
}

// Server 获取自己注入的组件服务
func (o *AppOptions) Server(name string) (any, bool) {
	srv, ok := o.servers[name]
	if !ok {
		return nil, false
	}
	return srv.Get(name), true
}

func WithAppID(id string) AppOption {
	return func(o *AppOptions) {
		if id != "" {
			o.appId = id
		}
	}
}

func WithSignal(sigs []os.Signal) AppOption {
	return func(o *AppOptions) {
		if len(sigs) > 0 {
			o.sigs = sigs
		}
	}
}

func WithEndpoints(endpoints []*url.URL) AppOption {
	return func(o *AppOptions) {
		if len(endpoints) > 0 {
			o.endpoints = endpoints
		}
	}
}

func WithRegistrar(registrar registry.ServiceRegistrar) AppOption {
	return func(o *AppOptions) {
		o.registrar = registrar
	}
}

func WithMustRedis() AppOption {
	components[compRedis] = struct{}{}
	return func(o *AppOptions) {
		o.wg.Go(func() error {
			var err error
			o.redis, err = cache.NewRedis(o.conf.Redis)
			if err != nil {
				return errors.Wrap(err, "init redis failed")
			}
			slog.Info(fmt.Sprintf(initDoneFmt, compRedis))
			return nil
		})
	}
}

func WithMustDB() AppOption {
	components[compDB] = struct{}{}
	return func(o *AppOptions) {
		o.wg.Go(func() error {
			smMap := make(map[string]db.InstanceConfig, len(o.conf.DB))
			for k, c := range o.conf.DB {
				smMap[k] = db.InstanceConfig{
					Driver: mysql.Open(c.Source), // TODO 支持其他数据库类型
					Config: c,
				}
			}
			var err error
			o.db, err = db.NewDBManager(smMap)
			if err != nil {
				return errors.Wrap(err, "init db manager failed")
			}
			slog.Info(fmt.Sprintf(initDoneFmt, compDB))
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

func WithGinServer(svr func(e *gin.Engine, a *AppOptions)) AppOption {
	return WithServer(compHttp, func(a *AppOptions) server.Server {
		return NewHttpServer(svr, a)
	})
}

func WithGrpcServer(svr func(s *grpc.Server, a *AppOptions), grpcServerOpts ...grpc.ServerOption) AppOption {
	return WithServer(compRpc, func(a *AppOptions) server.Server {
		return NewRpcServer(svr, a, grpcServerOpts...)
	})
}

func WithBeforeStart(fn func(ctx context.Context) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.beforeStart = append(o.beforeStart, fn)
		}
	}
}

func WithAfterStart(fn func(ctx context.Context, opts *AppOptions) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.afterStart = append(o.afterStart, fn)
		}
	}
}

func WithBeforeStop(fn func(ctx context.Context, opts *AppOptions) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.beforeStop = append(o.beforeStop, fn)
		}
	}
}

func WithAfterStop(fn func(ctx context.Context, opts *AppOptions) error) AppOption {
	return func(o *AppOptions) {
		if fn != nil {
			o.afterStop = append(o.afterStop, fn)
		}
	}
}
