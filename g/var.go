package g

import (
	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/db"
	"github.com/redis/go-redis/v9"
)

var (
	// Conf All 配置
	Conf *conf.Basic

	// CacheLocal
	// 1.一级缓存 变动小、容量少。容量固定，有淘汰策略。
	// 2.不适合分布式数据共享。
	CacheLocal cache.Cache

	// CacheDB 二级缓存 容量大，有网络IO延迟
	CacheDB redis.UniversalClient

	// DB gorm 关系型数据库 -- 持久化
	DB *db.DBManager
)
