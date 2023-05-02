package app

import (
	"context"
	"fmt"
	"github.com/gogoclouds/gogo/web/gin/valid"
	"time"

	"github.com/gogoclouds/gogo/internal/db"
	"github.com/patrickmn/go-cache"

	logger "github.com/gogoclouds/gogo/internal/log"
	"github.com/gogoclouds/gogo/internal/server"
	"github.com/gogoclouds/gogo/pkg/util"

	"github.com/gogoclouds/gogo/g"

	"github.com/gogoclouds/gogo/internal/conf"
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
	valid.InitRequestParamValidate()
	return &app{ctx: ctx, conf: g.Conf}
}

// Database connect DB
func (s *app) Database() *app {
	var err error
	if g.DB, err = db.Server.NewDB(s.ctx, s.conf); err != nil {
		panic(err)
	}
	return s
}

// AutoMigrate create data table
// tableModel struct 数据库表
func (s *app) AutoMigrate(tableModel []any) *app {
	if err := db.Server.AutoMigrate(g.DB, tableModel); err != nil {
		panic(err)
	}
	return s
}

func (s *app) Cache() *app {
	var err error
	if g.CacheDB, err = db.Redis.Open(s.ctx, s.conf); err != nil {
		panic(err)
	}
	return s
}

func (s *app) HTTP(router server.RegisterHttpFn) *app {
	httpConf := s.conf.App().Server.Http
	s.enableHttp = true
	go server.RunHttpServer(httpConf.Addr, router)
	return s
}

func (s *app) RPC(router server.RegisterRpcFn) *app {
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