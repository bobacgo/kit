package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type AppInfo struct {
	Name    string
	ID      string
	Version string
}

type TraceConfig struct {
	GrpcEndpoint string `mapstructure:"grpcEndpoint" yaml:"grpcEndpoint"`
}

type TracerServer struct {
	appInfo        AppInfo
	conf           *TraceConfig
	tracerProvider *sdktrace.TracerProvider
}

func TraceWrap(ctx context.Context, traceName, spanName string, fn func(ctx context.Context)) {
	ctx, span := otel.Tracer(traceName).Start(ctx, spanName)
	defer span.End()
	fn(ctx)
}

func NewTracerServer(appInfo AppInfo, conf *TraceConfig) *TracerServer {
	return &TracerServer{
		appInfo: appInfo,
		conf:    conf,
	}
}

func (srv *TracerServer) Start(ctx context.Context) error {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(srv.conf.GrpcEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("otlptracegrpc.New error %w: ", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(srv.appInfo.Name),
			semconv.ServiceInstanceID(srv.appInfo.ID),
			semconv.ServiceVersion(srv.appInfo.Version),
		),
	)
	if err != nil {
		return fmt.Errorf("resource.New error %w: ", err)
	}

	srv.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(srv.tracerProvider)
	return nil
}

func (srv *TracerServer) Stop(ctx context.Context) error {
	return srv.tracerProvider.Shutdown(context.Background())
}

func (srv *TracerServer) Get() any {
	return srv.tracerProvider
}
