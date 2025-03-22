package security

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bobacgo/kit/app/types"
	"github.com/bobacgo/kit/pkg/uid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const (
	CacheKeyPrefix                       = "login_token:"
	ATokenExpiredDuration types.Duration = "2h"
	RTokenExpiredDuration types.Duration = "360h" // 15天
)

type JWToken struct {
	cfg   JwtConfig
	cache redis.Cmdable
}

// TODO 实现本地存储
// func NewJWTLocal(conf *JwtConfig) *JWToken {
// 	jwt := &JWToken{cfg: *conf}
// 	jwt.init()
// 	return jwt
// }

func NewJWT(conf *JwtConfig, rdb redis.Cmdable) *JWToken {
	if conf.CacheKeyPrefix == "" {
		conf.CacheKeyPrefix = CacheKeyPrefix
	}
	jwt := &JWToken{cfg: *conf, cache: rdb}
	jwt.init()
	return jwt
}

func (t *JWToken) init() {
	if t.cfg.AccessTokenExpired == "" {
		t.cfg.AccessTokenExpired = ATokenExpiredDuration
	}
	if t.cfg.RefreshTokenExpired == "" {
		t.cfg.RefreshTokenExpired = RTokenExpiredDuration
	}
}

// Generate 颁发token access token 和 refresh token
// refresh token 不需要保存任何用户信息
func (t *JWToken) Generate(ctx context.Context, claims *Claims) (atoken, rtoken string, err error) {
	claims.ID = uid.UUID()
	claims.Issuer = t.cfg.Issuer
	claims.Audience = t.cfg.Audience
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(t.cfg.AccessTokenExpired.TimeDuration()))
	claims.NotBefore = jwt.NewNumericDate(time.Now().UTC())

	atoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.cfg.Secret)
	if err != nil {
		err = fmt.Errorf("access token generate err: %w", err)
		return
	}

	// refresh token 不需要保存任何用户信息
	sampleClains := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.cfg.RefreshTokenExpired.TimeDuration())),
		NotBefore: claims.NotBefore,
		Subject:   claims.Subject,
	}
	rtoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, sampleClains).SignedString(t.cfg.Secret)
	if err != nil {
		err = fmt.Errorf("refresh token generate err: %w", err)
		return
	}

	err = t.cacheToken(ctx, claims.Subject, claims.ID, atoken)
	return
}

func (t *JWToken) cacheToken(ctx context.Context, subject, tokenID, token string) error {
	value := fmt.Sprintf("%s|%s", tokenID, token)
	if t.cache == nil {
		return nil
	}
	return t.cache.Set(ctx, t.key(subject), value, t.cfg.AccessTokenExpired.TimeDuration()).Err()
}

func (t *JWToken) keyfunc(_ *jwt.Token) (any, error) {
	return t.cfg.Secret, nil
}

func (t *JWToken) Parse(tokenString string) (*Claims, error) {
	claims := new(Claims)
	_, err := jwt.ParseWithClaims(tokenString, claims, t.keyfunc)
	return claims, err
}

// Refresh 通过 refresh token 刷新 atoken
func (t *JWToken) Refresh(ctx context.Context, rtoken string, claims *Claims) (newAToken, newRToken string, err error) {
	if _, err = jwt.Parse(rtoken, t.keyfunc); err != nil { // rtoken 无效直接返回
		return
	}
	return t.Generate(ctx, claims)
}

// GetToken 获取 token
func (t *JWToken) GetToken(ctx context.Context, subject string) (string, error) {
	tokenInfo, err := t.getToken(ctx, subject)
	if err != nil {
		return "", err
	}
	if len(tokenInfo) < 2 {
		return "", errors.New("token extract fail")
	}
	return tokenInfo[1], nil
}

// GetTokenID 获取 tokenID
func (t *JWToken) GetTokenID(ctx context.Context, subject string) (string, error) {
	tokenInfo, err := t.getToken(ctx, subject)
	if err != nil {
		return "", err
	}
	if len(tokenInfo) < 2 {
		return "", errors.New("tokenID extract fail")
	}
	return tokenInfo[0], nil
}

// RemoveToken 删除 token
func (t *JWToken) RemoveToken(ctx context.Context, subject string) error {
	if t.cache == nil {
		return nil
	}
	return t.cache.Del(ctx, t.key(subject)).Err()
}

func (t *JWToken) getToken(ctx context.Context, subject string) ([]string, error) {
	if t.cache == nil {
		return nil, errors.New("token not cache")
	}
	tokenStr, err := t.cache.Get(ctx, t.key(subject)).Result()
	if err != nil {
		return nil, err
	}
	// 拆分 tokenID 和 token (tokenID|token)
	return strings.SplitN(tokenStr, "|", 2), nil
}

func (t *JWToken) key(subject string) string {
	return fmt.Sprintf("%s:%s", t.cfg.CacheKeyPrefix, subject)
}
