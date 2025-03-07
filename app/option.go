package app

import (
	"encoding/json"
	"github.com/bobacgo/kit/app/server"
	"github.com/bobacgo/kit/pkg/tag"
	"gopkg.in/yaml.v2"
	"log/slog"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/sync/errgroup"

	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/db"
	"github.com/bobacgo/kit/app/logger"
	"github.com/bobacgo/kit/app/registry"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
)

const (
	compLogger = "logger"
	compCache  = "local_cache"
	compRedis  = "redis"
	compDB     = "db"
	compHttp   = "http"
	compRpc    = "rpc"
)

var components = make(map[string]struct{}) // 组件名列表

type Option func(o *Options)

type Options struct {
	appId string      // 应用程序启动实例ID
	sigs  []os.Signal // 监听的程序退出信号

	conf *conf.App

	wgInit *errgroup.Group // 用于并发初始化组件
	// 内置功能
	localCache cache.Cache
	redis      redis.UniversalClient
	db         *db.DBManager
	// 插件功能 如 服务需要依赖 MongoDB、Elasticsearch等
	servers map[string]server.Server

	endpoints []*url.URL
	registrar registry.ServiceRegistrar

	httpServer func(e *gin.Engine, a *Options)  // 基于 gin 的 http 服务
	rpcServer  func(s *grpc.Server, a *Options) // 基于 grpc 的 rpc 服务
}

// AppID 获取应用程序启动实例ID
func (o *Options) AppID() string {
	return o.appId
}

// Conf 获取公共配置(eg app info、logger config、db config 、redis config)
func (o *Options) Conf() *conf.Basic {
	return &o.conf.Basic
}

// LocalCache 获取本地缓存 Interface
// CacheLocal
// 1.一级缓存 变动小、容量少。容量固定，有淘汰策略。
// 2.不适合分布式数据共享。
func (o *Options) LocalCache() cache.Cache {
	return o.localCache
}

// DB 获取数据库连接
// DB gorm 关系型数据库 -- 持久化
func (o *Options) DB() *db.DBManager {
	return o.db
}

// Redis 获取redis client
// CacheDB 二级缓存 容量大，有网络IO延迟
func (o *Options) Redis() redis.UniversalClient {
	return o.redis
}

// Server 获取自己注入的组件服务
func (o *Options) Server(name string) (any, bool) {
	srv, ok := o.servers[name]
	if !ok {
		return nil, false
	}
	return srv.Get(name), true
}

func WithAppID(id string) Option {
	return func(o *Options) {
		o.appId = id
	}
}

func WithSignal(sigs []os.Signal) Option {
	return func(o *Options) {
		o.sigs = sigs
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *Options) {
		o.endpoints = endpoints
	}
}

func WithRegistrar(registrar registry.ServiceRegistrar) Option {
	return func(o *Options) {
		o.registrar = registrar
	}
}

func WithScanConfig[T any](c *T) Option {
	return func(o *Options) {
		bytes, _ := json.Marshal(o.conf.Service)
		if err := json.Unmarshal(bytes, c); err != nil {
			slog.Error("scan config", "err", err)
		}
	}
}

func WithLogger() Option {
	components[compLogger] = struct{}{}
	return func(o *Options) {
		o.conf.Logger = logger.NewConfig()
		o.conf.Logger.Filename = o.conf.Name
		logger.InitZapLogger(o.conf.Logger)

		// 提供一个脱敏标签(mask)的配置文件
		// 扫描标签 并用 ** 替换
		maskConf := tag.Desensitize(o.conf)
		cfgData, _ := yaml.Marshal(maskConf)
		slog.Info("local config info\n" + string(cfgData))

		slog.Info("[config] init done.")
		slog.Info("[logger] init done.")
	}
}

// WithLocalCache 本地缓存
// 如果没有配置 LocalCache.MaxSize 则使用默认值 512MB
func WithLocalCache() Option {
	components[compCache] = struct{}{}
	return func(o *Options) {
		o.wgInit.Go(func() error {
			maxSize := o.conf.LocalCache.MaxSize
			if maxSize == "" {
				o.localCache = cache.DefaultCache()
			} else {
				var err error
				if o.localCache, err = cache.NewFreeCache(maxSize); err != nil {
					return errors.Wrap(err, "init local cache failed")
				}
			}
			slog.Info("[local_cache] init done.")
			return nil
		})
	}
}

func WithMustRedis() Option {
	components[compRedis] = struct{}{}
	return func(o *Options) {
		o.wgInit.Go(func() error {
			var err error
			o.redis, err = cache.NewRedis(o.conf.Redis)
			if err != nil {
				return errors.Wrap(err, "init redis failed")
			}
			slog.Info("[redis] init done.")
			return nil
		})
	}
}

func WithMustDB() Option {
	components[compDB] = struct{}{}
	return func(o *Options) {
		o.wgInit.Go(func() error {
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
			slog.Info("[database] init done.")
			return nil
		})
	}
}

func WithGinServer(svr func(e *gin.Engine, a *Options)) Option {
	components[compHttp] = struct{}{}
	return func(o *Options) {
		o.httpServer = svr
	}
}

func WithGrpcServer(svr func(s *grpc.Server, a *Options)) Option {
	components[compRpc] = struct{}{}
	return func(o *Options) {
		o.rpcServer = svr
	}
}

func WithServer(name string, srv func(a *Options) server.Server) Option {
	components[name] = struct{}{}
	return func(o *Options) {
		o.servers[name] = srv(o)
	}
}