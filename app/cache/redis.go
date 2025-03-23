package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// NewRedis Initialize redis connection.
func NewRedis(appName string, cfg RedisConf) (redis.UniversalClient, error) {
	if len(cfg.Addrs) == 0 {
		return nil, errors.New("redis address is empty")
	}
	var (
		err error
		rdb redis.UniversalClient
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if len(cfg.Addrs) > 1 {
		clt := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          cfg.Addrs, // []string{"IP_ADDRESS:6379"},
			ClientName:     appName,   // client name 方便监控和管理客户端连接
			ReadOnly:       true,      // 允许从节点可以执行读命令 GET、MGET、HGETALL、ZRANGE、SCAN 等，减轻主节点压力
			RouteByLatency: true,      // 按照延迟最低的节点进行读操作（主/从均可） 自动启用 ReadOnly
			Username:       cfg.Username,
			Password:       cfg.Password,
			PoolSize:       cfg.PoolSize,
			ReadTimeout:    cfg.ReadTimeout.TimeDuration(),  // 读取超时时间 默认是 3 秒
			WriteTimeout:   cfg.WriteTimeout.TimeDuration(), // 写入超时时间 默认是 5 秒
		})
		err = clt.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
			return shard.Ping(ctx).Err()
		})
		rdb = clt
	} else {
		clt := redis.NewClient(&redis.Options{
			Addr:         cfg.Addrs[0], // "IP_ADDRESS:6379",
			ClientName:   appName,      // client name 方便监控和管理客户端连接
			Username:     cfg.Username,
			Password:     cfg.Password,                    // no password set
			DB:           int(cfg.DB),                     // use default DB
			PoolSize:     cfg.PoolSize,                    // connection pool size 100
			ReadTimeout:  cfg.ReadTimeout.TimeDuration(),  // 读取超时时间 默认是 3 秒
			WriteTimeout: cfg.WriteTimeout.TimeDuration(), // 写入超时时间 默认是 5 秒
		})
		err = clt.Ping(ctx).Err()
		rdb = clt
	}
	if err != nil {
		return nil, fmt.Errorf("redis ping %v", err)
	}
	return rdb, nil
}
