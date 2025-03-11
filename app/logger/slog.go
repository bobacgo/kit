package logger

import (
	"log/slog"
)

// func InitSlog(handlers ...slog.Handler) {
// 	sc := slogcontext.NewHandler(slogmulti.Fanout(handlers...))
// 	l := slog.New(sc)
// 	slog.SetDefault(l)
// }

func InitSlog(h slog.Handler) {
	// sc := slogcontext.NewHandler(slogmulti.Fanout(handlers...))
	l := slog.New(h)
	slog.SetDefault(l)
}