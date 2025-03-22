package security

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	// https://tools.ietf.org/html/rfc7519  RFC 7519 定义的标准
	/*
		type StandardClaims struct {
			Audience  string `json:"aud,omitempty"`  // 受众（Audience），即该 JWT 令牌的目标用户或系统 (API 服务器的标识, https://api.example.com)
			ExpiresAt int64  `json:"exp,omitempty"`  // 过期时间（Expiration Time），以 UNIX 时间戳（秒）表示
			Id        string `json:"jti,omitempty"`  // JWT 唯一标识（JWT ID），用于避免令牌重放攻击
			IssuedAt  int64  `json:"iat,omitempty"`  // 签发时间（Issued At），表示令牌的创建时间
			Issuer    string `json:"iss,omitempty"`  // 签发者（Issuer），通常是颁发 JWT 的服务
			NotBefore int64  `json:"nbf,omitempty"`  // 生效时间（Not Before），表示该令牌在此时间之后才有效
			Subject   string `json:"sub,omitempty"`  // 主题（Subject），通常是用户 ID 或用户名
		}
	*/
	jwt.RegisteredClaims
	Data any `json:"data,omitempty"` // 自定义数据
}

const (
	ClaimsKey = "claims"
)

func GetUserInfo[T any](ctx context.Context) T {
	info := new(T)
	info, _ = ctx.Value(ClaimsKey).(*T)
	return *info
}
