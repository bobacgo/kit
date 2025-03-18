package client

import (
	"github.com/bobacgo/kit/app/conf"
	"github.com/bobacgo/kit/app/server/rpc/interceptor"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPC(transport conf.Transport, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if transport.Timeout == "" {
		transport.Timeout = "5s"
	}
	defaultOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 负载均衡策略 默认是 pick_first，所以我们换成 round_robin
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(transport.Timeout.TimeDuration()),
			logging.UnaryClientInterceptor(interceptor.Logger(), logging.WithFieldsFromContext(interceptor.LogTraceID))),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(interceptor.Logger(), logging.WithFieldsFromContext(interceptor.LogTraceID))),
	}

	defaultOpts = append(defaultOpts, opts...)
	return grpc.NewClient(transport.Addr, defaultOpts...)
}
