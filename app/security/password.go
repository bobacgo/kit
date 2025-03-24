package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bobacgo/kit/app/cache"
	"github.com/bobacgo/kit/pkg/ucrypto"
	"github.com/bobacgo/kit/pkg/utime"
	"github.com/redis/go-redis/v9"
)

var (
	ErrPasswdLimit = errors.New("password error limit")
)

// PasswdVerifier 登录密码验证器
// 1.对密码进行hash加密
// 2.随机生成盐
// 3.密码错误次数限制(依赖Redis)
type PasswdVerifier struct {
	cache      cache.Cache
	rdb        redis.Cmdable
	expiration time.Duration // 限制时长(在有效的错误次数范围内,每次错误都会刷新)
	limit      int32         // 错误次数限制
}

// DefaultPasswdVerifier 本地统计错误次数 (单节点)
func DefaultPasswdVerifier(cache cache.Cache, expiration time.Duration, limit int32) *PasswdVerifier {
	return &PasswdVerifier{
		cache:      cache,
		expiration: expiration,
		limit:      5,
	}
}

// NewPasswdVerifier 通过redis实现密码错误次数限制 (多节点)
// 1. keyTmp: 错误次数存放的key的模板 key = fmt.Sprintf(keyTmp, username)
// 2. 如果 expiration 为0,则使用默认的过期时间为第二天零点
func NewPasswdVerifier(rdb redis.Cmdable, expiration time.Duration, limit int32) *PasswdVerifier {
	return &PasswdVerifier{
		rdb:        rdb,
		expiration: expiration,
		limit:      limit,
	}
}

// BcryptVerify 验证密码
func (h *PasswdVerifier) BcryptVerify(hash, password string) bool {
	return verifyPwd(hash, password)
}

// BcryptHash 密码加密
func (h *PasswdVerifier) BcryptHash(passwd string) string {
	return processPwd(passwd)
}

// VerifierAndCount 验证密码统计错误次数
func (h *PasswdVerifier) VerifierAndCount(key string) PwdVerifier {
	return PwdVerifier{
		key: key,
		pv:  h,
	}
}

type PwdVerifier struct {
	key      string
	errCount int32
	pv       *PasswdVerifier
	OnErr    func(err error)
}

// BcryptVerifyWithCount 验证密码统计错误次数
func (h *PwdVerifier) BcryptVerify(ctx context.Context, hash, password string) bool {
	if len(hash) <= 8 {
		if h.OnErr != nil {
			h.OnErr(errors.New("hash length error"))
		}
		return false
	}
	if !verifyPwd(hash, password) {
		h.fail(ctx)
		return false
	}
	// 验证成功,删除错误次数
	if err := h.delIncr(ctx); err != nil {
		h.OnErr(err)
	}
	return true
}

// Incr 密码错误次数+1
func (h *PwdVerifier) Incr(ctx context.Context) {
	h.fail(ctx)
}

// Clear 清除密码错误次数
func (h *PwdVerifier) Clear(ctx context.Context) {
	if err := h.delIncr(ctx); err != nil && h.OnErr != nil {
		h.OnErr(err)
	}
	h.errCount = 0
}

// GetErrCount 获取密码错误次数
func (h *PwdVerifier) GetErrCount() int32 {
	return h.errCount
}

// GetRemainCount 获取密码剩余的错误次数
func (h *PwdVerifier) GetRemainCount() int32 {
	return max(h.pv.limit-int32(h.errCount), 0)
}

func (h *PwdVerifier) expire() time.Duration {
	if h.pv.expiration != 0 {
		return h.pv.expiration
	}
	return time.Duration(utime.ZeroHour(1).Unix() - time.Now().Unix())
}

func (h *PwdVerifier) fail(ctx context.Context) {
	if err := h.incr(ctx); err != nil && h.OnErr != nil {
		h.OnErr(err)
	}
	if err := h.reExpire(ctx); err != nil && h.OnErr != nil {
		h.OnErr(err)
	}
}

func (h *PwdVerifier) incr(ctx context.Context) error {
	var err error
	if h.pv.rdb != nil {
		count, ierr := h.pv.rdb.Incr(ctx, h.key).Result()
		if ierr != nil && !errors.Is(ierr, redis.Nil) {
			err = fmt.Errorf("redis incr %s: %w", h.key, ierr)
		} else {
			h.errCount = int32(count)
		}
	} else if h.pv.cache != nil {
		h.pv.cache.Get(h.key, &h.errCount)
		h.errCount++
	}

	if h.errCount >= h.pv.limit {
		return ErrPasswdLimit
	}
	return err
}

func (h *PwdVerifier) delIncr(ctx context.Context) error {
	if h.pv.rdb != nil {
		if err := h.pv.rdb.Del(ctx, h.key).Err(); err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis del %s: %w", h.key, err)
		}
	} else if h.pv.cache != nil {
		h.pv.cache.Del(h.key)
	}
	return nil
}

// 重置Key的过期时间
func (h *PwdVerifier) reExpire(ctx context.Context) error {
	if h.pv.rdb != nil {
		if err := h.pv.rdb.Expire(ctx, h.key, h.expire()).Err(); err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis expire %s: %w", h.key, err)
		}
	} else if h.pv.cache != nil {
		h.pv.cache.Set(h.key, h.errCount, h.expire())
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
