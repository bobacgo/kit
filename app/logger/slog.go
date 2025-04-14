package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// traceHandler 是一个添加 trace ID 和 span ID 的 slog handler
// traceHandler is a slog handler that adds trace ID and span ID
type traceHandler struct {
	handler slog.Handler
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, record slog.Record) error {
	// 从 context 中提取 trace 信息
	// Extract trace information from context
	if ctx != nil {
		spanCtx := trace.SpanContextFromContext(ctx)
		if spanCtx.IsValid() {
			record.AddAttrs(
				slog.String("trace_id", spanCtx.TraceID().String()),
				slog.String("span_id", spanCtx.SpanID().String()),
			)
		}
	}
	return h.handler.Handle(ctx, record)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{
		handler: h.handler.WithGroup(name),
	}
}

func InitSlog(h slog.Handler) {
	// 包装原始 handler，添加 trace 支持
	// Wrap the original handler with trace support
	tracingHandler := &traceHandler{handler: h}
	l := slog.New(tracingHandler)
	slog.SetDefault(l)
}
