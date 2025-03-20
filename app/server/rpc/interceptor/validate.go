package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bobacgo/kit/app/validator"
)

// ValidateInterceptor 请求参数校验拦截器
func ValidateParam() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		lang := GetLanguage(ctx)
		// 对请求参数进行校验
		if err := validator.StructLocale(lang, req); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		// 继续处理请求
		return handler(ctx, req)
	}
}

// ValidateStreamInterceptor 流式请求参数校验拦截器
func ValidateStreamParam() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 包装ServerStream以拦截RecvMsg和SendMsg
		wrappedStream := &validatorServerStream{
			ServerStream: ss,
			ctx:          ss.Context(),
		}
		return handler(srv, wrappedStream)
	}
}

type validatorServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *validatorServerStream) Context() context.Context {
	return s.ctx
}

func (s *validatorServerStream) RecvMsg(m any) error {
	// 先接收消息
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	lang := GetLanguage(s.Context())

	// 对接收到的消息进行校验
	if err := validator.StructLocale(lang, m); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func (s *validatorServerStream) SendMsg(m any) error {
	return s.ServerStream.SendMsg(m)
}
