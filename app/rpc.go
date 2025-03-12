package app

import (
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type RpcServer struct {
	Opts       *Options
	RegistryFn func(s *grpc.Server, a *Options)

	server *grpc.Server
}

func NewRpcServer(registry func(s *grpc.Server, a *Options), opts *Options) *RpcServer {
	return &RpcServer{
		Opts:       opts,
		RegistryFn: registry,
	}
}

func (srv *RpcServer) Get(name string) any {
	return nil
}

func (srv *RpcServer) Start(ctx context.Context) error {
	cfg := srv.Opts.Conf()

	srv.server = grpc.NewServer()
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
