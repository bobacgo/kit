package interceptor

import (
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(p any) error {
	slog.Error("recovered from panic", "panic", p, "stack", debug.Stack())
	return status.Errorf(codes.Internal, "%s", p)
}
