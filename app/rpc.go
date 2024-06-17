package app

import (
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func RunMustRpcServer(app *App, register func(server *grpc.Server, app *Options)) {
	app.wg.Add(1)
	defer app.wg.Done()

	cfg := app.opts.Conf()

	lis, err := net.Listen("tcp", cfg.Server.Rpc.Addr)
	if err != nil {
		log.Panic(err)
	}
	s := grpc.NewServer()

	// 注册健康检查服务
	healthgrpc.RegisterHealthServer(s, health.NewServer())

	if register != nil {
		register(s, &app.opts)
	}

	go func() {
		localhost, _ := getRegistryUrl("grpc", cfg.Server.Http.Addr)
		slog.Info("grpc server running " + localhost)

		if err = s.Serve(lis); err != nil {
			log.Panic(err)
		}
	}()

	<-app.exit       // 阻塞,等待被关闭
	s.GracefulStop() // 优雅停止
}
