package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bobacgo/kit/app/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GatewayRegisterFunc defines the function type for registering handlers to the gateway
// 定义网关处理器注册函数类型
type GatewayRegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

// GatewayRegisterItem defines the registration function and its corresponding endpoint
// 定义注册项，包含注册函数及其对应的端点
type GatewayRegisterItem struct {
	Func     GatewayRegisterFunc // Registration function // 注册函数
	Endpoint string              // gRPC endpoint address // gRPC 端点地址
	DialOpts []grpc.DialOption   // gRPC dial options // gRPC 拨号选项
}

// DefaultGatewayConfig returns a default configuration for gateway
// 返回默认网关配置
func DefaultGatewayConfig() *Config {
	return &Config{
		Addr:            ":8080",
		ReadTimeout:     "10s",
		WriteTimeout:    "10s",
		ShutdownTimeout: "15s",
	}
}

// Gateway implements the Server interface for grpc-gateway
// Gateway 实现了 grpc-gateway 的 Server 接口
type Gateway struct {
	cfg           *Config
	RegisterItems []GatewayRegisterItem // Registration items with endpoints // 带端点的注册项
	server        *http.Server
	mux           *runtime.ServeMux
}

// New creates a new Gateway instance
// 创建一个新的 Gateway 实例
func New(cfg *Config, register []GatewayRegisterItem, serveMuxOptions ...runtime.ServeMuxOption) *Gateway {
	if cfg == nil {
		cfg = DefaultGatewayConfig()
	}

	serveMuxOpts := []runtime.ServeMuxOption{
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	}
	serveMuxOpts = append(serveMuxOpts, serveMuxOptions...)
	mux := runtime.NewServeMux(serveMuxOpts...)
	return &Gateway{
		cfg:           cfg,
		RegisterItems: register,
		mux:           mux,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout.TimeDuration(),
			WriteTimeout: cfg.WriteTimeout.TimeDuration(),
		},
	}
}

// ServeSwaggerUI serves swagger UI if enabled
// 如果启用了swagger，则提供swagger UI服务
func (g *Gateway) ServeSwaggerUI(w http.ResponseWriter, r *http.Request, pathPrefix string) {
	if g.cfg.SwaggerDir == "" {
		http.NotFound(w, r)
		return
	}

	// Redirect to swagger UI
	// 重定向到swagger UI
	if r.URL.Path == pathPrefix {
		http.Redirect(w, r, pathPrefix+"/", http.StatusMovedPermanently)
		return
	}

	// Serve swagger-ui and swagger.json
	// 提供swagger-ui和swagger.json
	if r.URL.Path == pathPrefix+"/" {
		http.ServeFile(w, r, g.cfg.SwaggerDir+"/swagger-ui/index.html")
	} else if strings.HasPrefix(r.URL.Path, pathPrefix+"/") {
		// Serve other swagger-ui resources
		// 提供其他swagger-ui资源
		http.StripPrefix(pathPrefix+"/", http.FileServer(http.Dir(g.cfg.SwaggerDir+"/swagger-ui"))).ServeHTTP(w, r)
	} else if strings.HasSuffix(r.URL.Path, ".swagger.json") {
		// Serve the specific swagger definition JSON
		// 提供特定的swagger定义JSON
		http.ServeFile(w, r, g.cfg.SwaggerDir+r.URL.Path)
		return
	}

	http.NotFound(w, r)
}

// Start implements the Server interface
// 实现 Server 接口的 Start 方法
func (g *Gateway) Start(ctx context.Context) error {
	// Register all handlers
	// 注册所有处理器
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
	}

	// Check if we have any registration items
	// 检查是否有注册项
	if len(g.RegisterItems) == 0 {
		return fmt.Errorf("no gRPC service registration items provided")
	}

	// Register each function with its corresponding endpoint
	// 为每个注册函数使用其对应的端点进行注册
	for _, item := range g.RegisterItems {
		dialOpts := append([]grpc.DialOption{}, opts...)
		dialOpts = append(dialOpts, item.DialOpts...)
		if err := item.Func(ctx, g.mux, item.Endpoint, dialOpts); err != nil {
			return fmt.Errorf("failed to register handler for endpoint %s: %w", item.Endpoint, err)
		}
	}

	// Create final handler with optional swagger support
	// 创建最终的处理器，可选择性地支持swagger
	var handler http.Handler = g.mux

	if g.cfg.SwaggerDir != "" {
		// Serve swagger
		// 提供swagger服务
		handler = g.wrapHandlerWithSwagger(handler)
	}

	// Update server handler
	// 更新服务器处理器
	g.server.Handler = handler

	// Start HTTP server
	// 启动 HTTP 服务器
	go func() {
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("gateway server failed to serve", "error", err)
		}
	}()

	// Log server started info
	// 记录服务器启动信息
	slog.Info("gateway server started", "addr", g.cfg.Addr)

	// Log swagger URL if enabled
	// 如果启用了swagger，记录swagger访问地址
	if g.cfg.SwaggerDir != "" {
		swaggerURL := fmt.Sprintf("http://%s/swagger/", g.server.Addr)
		if g.server.Addr[0] == ':' {
			// If the address starts with ':', it's a port number
			// 如果地址以':'开头，则为端口号
			swaggerURL = fmt.Sprintf("http://localhost%s/swagger/", g.server.Addr)
		}
		slog.Info("swagger UI available", "url", swaggerURL)
	}

	return nil
}

// wrapHandlerWithSwagger adds swagger handlers to the provided handler
// 为提供的处理器添加swagger处理器
func (g *Gateway) wrapHandlerWithSwagger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve swagger-ui and swagger.json
		// 提供swagger-ui和swagger.json
		if strings.HasPrefix(r.URL.Path, "/swagger") {
			if r.URL.Path == "/swagger" || r.URL.Path == "/swagger/" {
				http.ServeFile(w, r, g.cfg.SwaggerDir+"/swagger-ui/index.html")
				return
			} else if strings.HasPrefix(r.URL.Path, "/swagger/") {
				// Serve other swagger-ui resources
				// 提供其他swagger-ui资源
				http.StripPrefix("/swagger/", http.FileServer(http.Dir(g.cfg.SwaggerDir+"/swagger-ui"))).ServeHTTP(w, r)
				return
			}
		} else if strings.HasSuffix(r.URL.Path, ".swagger.json") {
			// Serve the specific swagger definition JSON
			// 提供特定的swagger定义JSON
			http.ServeFile(w, r, g.cfg.SwaggerDir+r.URL.Path)
			return
		}

		// Fall back to the original handler
		// 回退到原始处理器
		h.ServeHTTP(w, r)
	})
}

// Stop implements the Server interface
// 实现 Server 接口的 Stop 方法
func (g *Gateway) Stop(ctx context.Context) error {
	// Create a context with timeout for shutdown
	// 创建一个带超时的上下文用于关闭
	shutdownCtx, cancel := context.WithTimeout(ctx, g.cfg.ShutdownTimeout.TimeDuration())
	defer cancel()

	// Gracefully shutdown the server
	// 优雅关闭服务器
	if err := g.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gateway server shutdown failed: %w", err)
	}
	return nil
}

// Get implements the Server interface
// 实现 Server 接口的 Get 方法
func (g *Gateway) Get() any {
	return g.server
}

// Ensure Gateway implements the Server interface
// 确保 Gateway 实现了 Server 接口
var _ server.Server = (*Gateway)(nil)
