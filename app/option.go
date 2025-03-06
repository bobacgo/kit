package app

import (
	"encoding/json"
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

type Option func(o *Options)

type Options struct {
	appid string      // 应用程序启动实例ID
	sigs  []os.Signal // 监听的程序退出信号

	conf *conf.App

	wg         *errgroup.Group
	localCache cache.Cache
	redis      redis.UniversalClient
	db         *db.DBManager

	endpoints []*url.URL
	registrar registry.ServiceRegistrar

	httpServer func(e *gin.Engine, a *Options)
	rpcServer  func(s *grpc.Server, a *Options)
}

// AppID 获取应用程序启动实例ID
func (o *Options) AppID() string {
	return o.appid
}

// Conf 获取公共配置(eg app info、logger config、db config 、redis config)
func (o *Options) Conf() *conf.Basic {
	return &o.conf.Basic
}

// LocalCache 获取本地缓存 Interface
// CacheLocal
// 1.一级缓存 变动小、容量少。容量固定，有淘汰策略。
// 2.不适合分布式数据共享。
func (o Options) LocalCache() cache.Cache {
	return o.localCache
}

// DB 获取数据库连接
// DB gorm 关系型数据库 -- 持久化
func (o Options) DB() *db.DBManager {
	return o.db
}

// Redis 获取redis client
// CacheDB 二级缓存 容量大，有网络IO延迟
func (o Options) Redis() redis.UniversalClient {
	return o.redis
}

func WithAppID(id string) Option {
	return func(o *Options) {
		o.appid = id
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
		bytes, err := json.Marshal(o.conf.Service)
		if err != nil {
			return
		}
		_ = json.Unmarshal(bytes, c)
	}
}

func WithLogger() Option {
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

func WithMustLocalCache() Option {
	return func(o *Options) {
		o.wg.Go(func() error {
			var err error
			o.localCache, err = cache.DefaultCache()
			if err != nil {
				return errors.Wrap(err, "init local cache failed")
			}
			slog.Info("[local_cache] init done.")
			return nil
		})
	}
}

func WithMustRedis() Option {
	return func(o *Options) {
		o.wg.Go(func() error {
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
	return func(o *Options) {
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
			slog.Info("[database] init done.")
			return nil
		})
	}
}

func WithGinServer(router func(e *gin.Engine, a *Options)) Option {
	return func(o *Options) {
		o.httpServer = router
	}
}

func WithGrpcServer(svr func(s *grpc.Server, a *Options)) Option {
	return func(o *Options) {
		o.rpcServer = svr
	}
}