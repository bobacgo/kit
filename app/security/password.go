package security

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/bobacgo/kit/pkg/ucrypto"
	"github.com/bobacgo/kit/pkg/utime"
	"github.com/redis/go-redis/v9"
)

var (
	ErrPasswdLimit    = errors.New("password error limit")
	ErrPasswdMismatch = errors.New("password mismatch")
)

// PasswdVerifier 登录密码验证器
// 1.对密码进行hash加密
// 2.随机生成盐
// 3.密码错误次数限制(依赖Redis)
type PasswdVerifier struct {
	cache      redis.Cmdable
	key        string        // 错误密码存放的key
	expiration time.Duration // 限制时长(在有效的错误次数范围内,每次错误都会刷新)
	limit      int32         // 错误次数限制
	errCount   atomic.Int32  // 尝试次数
}

// DefaultPasswdVerifier 本地统计错误次数 (单节点)
func DefaultPasswdVerifier(limit int32) *PasswdVerifier {
	return &PasswdVerifier{limit: 5}
}

// NewPasswdVerifier 通过redis实现密码错误次数限制 (多节点)
// 1. 如果 expiration 为0,则使用默认的过期时间为第二天零点
func NewPasswdVerifier(rdb redis.Cmdable, key string, expiration time.Duration, limit int32) *PasswdVerifier {
	return &PasswdVerifier{
		cache:      rdb,
		key:        key,
		expiration: expiration,
		limit:      limit,
	}
}

func (h *PasswdVerifier) expire() time.Duration {
	if h.expiration != 0 {
		return h.expiration
	}
	return time.Duration(utime.ZeroHour(1).Unix() - time.Now().Unix())
}

// BcryptVerifyWithCount 验证密码统计错误次数
func (h *PasswdVerifier) BcryptVerifyWithCount(ctx context.Context, hash, password string) (bool, error) {
	if len(hash) <= 8 {
		return false, errors.New("hash length error")
	}
	if !verifyPwd(hash, password) {
		if err := h.fail(ctx); err != nil {
			return false, err
		}
		return false, ErrPasswdMismatch
	}

	// 验证成功,删除错误次数
	return true, h.delIncr(ctx)
}

// BcryptVerify 验证密码
func (h *PasswdVerifier) BcryptVerify(hash, password string) bool {
	return verifyPwd(hash, password)
}

// BcryptHash 密码加密
func (h *PasswdVerifier) BcryptHash(passwd string) string {
	return processPwd(passwd)
}

// GetErrCount 获取密码错误的次数
func (h *PasswdVerifier) GetErrCount() int32 {
	return h.errCount.Load()
}

// GetRemainCount 获取密码剩余的错误次数
func (h *PasswdVerifier) GetRemainCount() int32 {
	return max(h.limit-h.errCount.Load(), 0)
}

// Incr 密码错误次数+1
func (h *PasswdVerifier) Incr(ctx context.Context) error {
	if h.cache != nil {
		count, err := h.cache.Incr(ctx, h.key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis incr %s: %w", h.key, err)
		}
		h.errCount.Store(int32(count))
	} else {
		h.errCount.Add(1)
	}
	return nil
}

func (h *PasswdVerifier) fail(ctx context.Context) error {
	if err := h.Incr(ctx); err != nil {
		return err
	}
	if h.cache != nil {
		// 重置过期时间
		if err := h.cache.Expire(ctx, h.key, h.expire()).Err(); err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis expire %s: %w", h.key, err)
		}
	}
	if h.errCount.Load() >= h.limit {
		return ErrPasswdLimit
	}
	return nil
}

func (h *PasswdVerifier) delIncr(ctx context.Context) error {
	if h.cache == nil {
		return nil
	}
	if err := h.cache.Del(ctx, h.key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("redis del %s: %w", h.key, err)
	}
	return nil
}

// processPwd 处理密码
// 返回 hash + salt 拼接在一起
func processPwd(passwd string) string {
	hash, salt := ucrypto.BcryptHash(passwd)
	return hash + salt
}

// 解析密码
// hash + salt 后面8位为salt
func verifyPwd(hashPwd, password string) bool {
	if len(hashPwd) <= 8 {
		return false
	}
	hash, salt := hashPwd[:len(hashPwd)-8], hashPwd[len(hashPwd)-8:]
	return ucrypto.BcryptVerify(hash, salt, password)
}
