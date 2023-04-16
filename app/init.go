package app

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gogoclouds/gogo/internal/db"
	"github.com/patrickmn/go-cache"

	logger "github.com/gogoclouds/gogo/internal/log"
	"github.com/gogoclouds/gogo/internal/server"
	"github.com/gogoclouds/gogo/pkg/util"

	"github.com/gogoclouds/gogo/g"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gogoclouds/gogo/internal/conf"
	"github.com/gogoclouds/gogo/web/gin/valid"
)

type app struct {
	ctx        context.Context
	conf       *conf.Config
	enableRpc  bool
	enableHttp bool
}

// New().OpenDB().OpenCacheDB().CreateXxxServer().Run()

// New 这个函数调用之后会阻塞
// 1. 从配置中心拉取配置文件
// 2. 启动服务
// 3. 初始必要的全局参数
func New(ctx context.Context, configPath string) *app {
	g.Conf = conf.New(configPath)
	g.CacheLocal = cache.New(5*time.Minute, 10*time.Minute)

	logger.Initialize(g.Conf.App().Name, g.Conf.Log())
	initValidate()
	return &app{ctx: ctx, conf: g.Conf}
}

// OpenDB connect DB
//
// tableModel struct 数据库表
func (s *app) OpenDB(tableModel []any) *app {
	var err error
	if g.DB, err = db.Server.NewDB(s.ctx, s.conf); err != nil {
		panic(err)
	}
	if err = db.Server.AutoMigrate(g.DB, tableModel); err != nil {
		panic(err)
	}
	return s
}

func (s *app) OpenCacheDB() *app {
	var err error
	if g.CacheDB, err = db.Redis.Open(s.ctx, s.conf); err != nil {
		panic(err)
	}
	return s
}

func (s *app) CreateHttpServer(router server.RegisterHttpFn) *app {
	httpConf := s.conf.App().Server.Http
	s.enableHttp = true
	go server.RunHttpServer(httpConf.Addr, router)
	return s
}

func (s *app) CreateRpcServer(router server.RegisterRpcFn) *app {
	rpcConf := s.conf.App().Server.Rpc
	s.enableRpc = true
	go server.RunRpcServer(rpcConf.Addr, router)
	return s
}

func (s *app) Run() {
	var port uint16
	if s.enableHttp {
		_, port = util.IP.Parse(s.conf.App().Server.Http.Addr)
	}
	ip, _ := util.IP.GetOutBoundIP()
	fmt.Printf("http://%s:%d/health\n", ip, port)
	select {}
}

// initTrans 初始化翻译器
func initValidate() {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// skip if tag key says it should be ignored
			if name == "-" {
				return ""
			}
			return name
		})

		uni := ut.New(zh.New())
		valid.Trans, _ = uni.GetTranslator("zh")
		if err := zhTranslations.RegisterDefaultTranslations(validate, valid.Trans); err != nil {
			panic(err)
		}
	}
}
