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
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(transport.Timeout.TimeDuration()),
			logging.UnaryClientInterceptor(interceptor.Logger(), logging.WithFieldsFromContext(interceptor.LogTraceID))),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(interceptor.Logger(), logging.WithFieldsFromContext(interceptor.LogTraceID))),
	)
	return grpc.NewClient(transport.Addr, opts...)
}
