package gateway

import (
	"github.com/bobacgo/kit/app/types"
)

// Config holds the configuration for gateway server
// 网关服务器配置
type Config struct {
	Addr            string         `mapstructure:"addr" yaml:"addr"`                       // HTTP server address // HTTP服务器地址
	SwaggerDir      string         `mapstructure:"swaggerDir" yaml:"swaggerDir"`           // Directory containing swagger files // 包含swagger文件的目录
	ReadTimeout     types.Duration `mapstructure:"readTimeout" yaml:"readTimeout"`         // HTTP read timeout // HTTP读取超时时间
	WriteTimeout    types.Duration `mapstructure:"writeTimeout" yaml:"writeTimeout"`       // HTTP write timeout // HTTP写入超时时间
	ShutdownTimeout types.Duration `mapstructure:"shutdownTimeout" yaml:"shutdownTimeout"` // Graceful shutdown timeout // 优雅关闭超时时间
}
