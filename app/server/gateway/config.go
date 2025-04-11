package gateway

import (
	"github.com/bobacgo/kit/app/types"
)

// Config holds the configuration for gateway server
// 网关服务器配置
type Config struct {
	Addr            string         // HTTP server address // HTTP服务器地址
	SwaggerDir      string         // Directory containing swagger files // 包含swagger文件的目录
	ReadTimeout     types.Duration // HTTP read timeout // HTTP读取超时时间
	WriteTimeout    types.Duration // HTTP write timeout // HTTP写入超时时间
	ShutdownTimeout types.Duration // Graceful shutdown timeout // 优雅关闭超时时间
}
