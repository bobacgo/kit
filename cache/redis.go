package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

const (
	maxRetry = 10 // number of retries
)

// NewRedis Initialize redis connection.
func NewRedis(cfg RedisConf) (redis.UniversalClient, error) {
	if len(cfg.Addrs) == 0 {
		return nil, errors.New("redis address is empty")
	}
	var rdb redis.UniversalClient

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if len(cfg.Addrs) > 1 {
		clt := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.Addrs,
			Username:     cfg.Username,
			Password:     cfg.Password, // no password set
			PoolSize:     cfg.PoolSize, // 50
			MaxRetries:   maxRetry,
			ReadTimeout:  cfg.ReadTimeout.TimeDuration(),
			WriteTimeout: cfg.WriteTimeout.TimeDuration(),
		})
		err = clt.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
			return shard.Ping(ctx).Err()
		})
		rdb = clt
	} else {
		clt := redis.NewClient(&redis.Options{
			Addr:         cfg.Addrs[0],
			Username:     cfg.Username,
			Password:     cfg.Password, // no password set
			DB:           int(cfg.DB),  // use default DB
			PoolSize:     cfg.PoolSize, // connection pool size 100
			MaxRetries:   maxRetry,
			ReadTimeout:  cfg.ReadTimeout.TimeDuration(),
			WriteTimeout: cfg.WriteTimeout.TimeDuration(),
		})
		err = clt.Ping(ctx).Err()
		rdb = clt
	}

	return rdb, fmt.Errorf("redis ping %w", err)
}
