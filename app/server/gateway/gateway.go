package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bobacgo/kit/app/server"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// RegisterHandler registers a handler function with gin
// 注册一个处理器函数与gin
// 支持使用 gin 实现 grpc-gateway 的 http server
func RegisterHandler(e *gin.Engine, handlerFunc func(mux *runtime.ServeMux), serveMuxOptions ...runtime.ServeMuxOption) {
	mux := runtime.NewServeMux(serveMuxOptions...)
	if handlerFunc != nil {
		handlerFunc(mux)
	}
	// Register the handler with gin
	// 使用gin注册处理器
	e.NoRoute(func(c *gin.Context) {
		// Serve the request using the mux
		// 使用mux处理请求
		mux.ServeHTTP(c.Writer, c.Request)
	})
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
	cfg         *Config
	server      *http.Server
	mux         *runtime.ServeMux
	handlerFunc func(mux *runtime.ServeMux)
}

// New creates a new Gateway instance
// 创建一个新的 Gateway 实例
func New(cfg *Config, handlerFunc func(mux *runtime.ServeMux), serveMuxOptions ...runtime.ServeMuxOption) *Gateway {
	if cfg == nil {
		cfg = DefaultGatewayConfig()
	}

	serveMuxOpts := []runtime.ServeMuxOption{
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	}

	serveMuxOpts = append(serveMuxOpts, serveMuxOptions...)
	mux := runtime.NewServeMux(serveMuxOpts...)
	return &Gateway{
		cfg:         cfg,
		mux:         mux,
		handlerFunc: handlerFunc,
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
	if g.handlerFunc != nil {
		// Call the handler function to register handlers
		// 调用处理器函数以注册处理器
		g.handlerFunc(g.mux)
	}

	// Create final handler with optional swagger support
	// 创建最终的处理器，可选择性地支持swagger
	var handler http.Handler = g.mux

	if g.cfg.SwaggerDir != "" {
		// Serve swagger
		// 提供swagger服务
		handler = g.wrapHandlerWithSwagger(handler)
	}

	// 自动追踪的 HTTP handler
	handler = otelhttp.NewHandler(handler, "grpc-gateway")

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
