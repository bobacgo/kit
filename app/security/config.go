package security

import (
	"github.com/bobacgo/kit/app/types"
)

type Config struct {
	Ciphertext CiphertextConfig `mapstructure:"ciphertext"`
	Jwt        JwtConfig        `mapstructure:"jwt"`
}

type CiphertextConfig struct {
	IsCiphertext bool       `mapstructure:"isCiphertext" yaml:"isCiphertext"`   // 密码字段是否启用密文传输
	CipherKey    Ciphertext `mapstructure:"cipherKey" yaml:"cipherKey" mask:""` // 支持 8 16 24 bit
}

type JwtConfig struct {
	Secret         Ciphertext `mapstructure:"secret" mask:""`
	CacheKeyPrefix string     `mapstructure:"cacheKeyPrefix" yaml:"cacheKeyPrefix"` // jwt cache key prefix 分布式共享token
	// Claims jwt claims
	Audience            []string       `mapstructure:"audience"`                                                                          // jwt audience
	Issuer              string         `mapstructure:"issuer"`                                                                            // jwt issuer
	AccessTokenExpired  types.Duration `mapstructure:"accessTokenExpired" yaml:"accessTokenExpired" validate:"duration" default:"2h"`     // jwt access token expired
	RefreshTokenExpired types.Duration `mapstructure:"refreshTokenExpired" yaml:"refreshTokenExpired" validate:"duration" default:"720h"` // jwt refresh token expired
}

// TODO validate config

func (c *JwtConfig) Validate() []error {
	// TODO valid config data
	return nil
}
