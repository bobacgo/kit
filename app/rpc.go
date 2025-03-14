package app

import (
	"context"
	"log"
	"log/slog"
	"net"

	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/bobacgo/kit/app/server/rpc/interceptor"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type RpcServer struct {
	Opts       *AppOptions
	RegistryFn func(s *grpc.Server, a *AppOptions)

	grpcServerOpts []grpc.ServerOption
	server         *grpc.Server
}

func NewRpcServer(registry func(s *grpc.Server, a *AppOptions), opts *AppOptions, grpcServerOpts ...grpc.ServerOption) *RpcServer {
	return &RpcServer{
		Opts:           opts,
		RegistryFn:     registry,
		grpcServerOpts: grpcServerOpts,
	}
}

func (srv *RpcServer) Get(name string) any {
	return nil
}

func (srv *RpcServer) Start(ctx context.Context) error {
	cfg := srv.Opts.Conf()

	srv.defaultInterceptor()
	srv.server = grpc.NewServer(srv.grpcServerOpts...)

	healthgrpc.RegisterHealthServer(srv.server, health.NewServer()) // 注册健康检查服务
	if srv.RegistryFn != nil {                                      // 注册业务接口
		srv.RegistryFn(srv.server, srv.Opts)
	}

	// 保证端口监听成功
	lis, err := net.Listen("tcp", cfg.Server.Rpc.Addr)
	if err != nil {
		return err
	}

	localhost, _ := getRegistryUrl("grpc", cfg.Server.Rpc.Addr)
	slog.Info("grpc server running " + localhost)
	go func(lis net.Listener) {
		if err := srv.server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Panicf("listen: %s\n", err)
		}
	}(lis)
	return nil
}

func (srv *RpcServer) Stop(ctx context.Context) error {
	if srv.server == nil {
		return nil
	}

	slog.Info("Shutting down grpc server...")
	srv.server.GracefulStop() // 优雅停止
	return nil
}

func (srv *RpcServer) defaultInterceptor() {
	srv.grpcServerOpts = append(srv.grpcServerOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor( // 单向拦截器
			logging.UnaryServerInterceptor(interceptor.Logger(), logging.WithFieldsFromContext(interceptor.LogTraceID)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(interceptor.Recovery)),
		), grpc.ChainStreamInterceptor( // 流式拦截器
			logging.StreamServerInterceptor(interceptor.Logger(), logging.WithFieldsFromContext(interceptor.LogTraceID)),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(interceptor.Recovery)),
		))
}
