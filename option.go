package kit

import (
	"log"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/bobacgo/kit/cache"
	"github.com/bobacgo/kit/conf"
	"github.com/bobacgo/kit/db"
	"github.com/bobacgo/kit/g"
	"github.com/bobacgo/kit/logger"
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

	conf       *conf.Basic
	localCache cache.Cache
	redis      redis.UniversalClient
	db         *db.SourceManager

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
func (o Options) DB() *db.SourceManager {
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
		var err error
		o.localCache, err = cache.DefaultCache()
		if err != nil {
			log.Panic(err)
		}
		slog.Info("[local_cache] init done.")
	}
}

func WithMustRedis() Option {
	return func(o *Options) {
		var err error
		o.redis, err = cache.NewRedis(o.conf.Redis)
		if err != nil {
			log.Panic(err.Error())
		}
		slog.Info("[redis] init done.")
	}
}

func WithMustDB() Option {
	return func(o *Options) {
		smMap := make(map[string]db.InstanceConfig, len(o.conf.DB))
		for k, c := range o.conf.DB {
			smMap[k] = db.InstanceConfig{
				Driver: mysql.Open(c.Source), // TODO 支持其他数据库类型
				Config: c,
			}
		}
		o.db = db.NewSourceManager(smMap)
		slog.Info("[mysql] init done.")
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