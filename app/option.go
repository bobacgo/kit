package app

import (
	"log"
	"log/slog"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"

	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/db"
	"github.com/bobacgo/kit/app/logger"
	"github.com/bobacgo/kit/app/registry"
	"github.com/bobacgo/kit/g"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
)

type Option func(o *Options)

type Options struct {
	appid string      // 应用程序启动实例ID
	sigs  []os.Signal // 监听的程序退出信号

	conf *conf.Basic

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
	return o.conf
}

// LocalCache 获取本地缓存 Interface
func (o Options) LocalCache() cache.Cache {
	return o.localCache
}

// DB 获取数据库连接
func (o Options) DB() *db.DBManager {
	return o.db
}

// Redis 获取redis client
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

func WithMustConfig[T any](filename string, fn func(cfg *conf.ServiceConfig[T])) Option {
	return func(o *Options) {
		cfg, err := conf.LoadService[T](filename, func(e fsnotify.Event) {
			//logger.S(config.Conf.Logger.Level)
		})
		if err != nil {
			log.Panic(err)
		}
		o.conf = &cfg.Basic
		g.Conf = &cfg.Basic
		cfgData, _ := yaml.Marshal(cfg)
		slog.Info("local config info\n" + string(cfgData)) // TODO 脱密
		fn(cfg)
	}
}

func WithLogger() Option {
	return func(o *Options) {
		o.conf.Logger = logger.NewConfig()
		o.conf.Logger.Filename = o.conf.Name
		logger.InitZapLogger(o.conf.Logger)

		// yaml 格式输出到控制台
		// cfgData, _ := yaml.Marshal(viper.AllSettings())
		// slog.Info("local config info\n" + string(cfgData)) // TODO 脱密
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
			g.CacheLocal = o.localCache
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
			g.CacheDB = o.redis
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
			g.DB = o.db
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
