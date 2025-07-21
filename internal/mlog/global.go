package mlog

import (
	"log/slog"
	"sync/atomic"
)

var Log func() *slog.Logger = func() *slog.Logger {
	logger := defaultLoggerP.Load()
	if logger == nil {
		return slog.Default()
	}
	return logger
}
var defaultLoggerP atomic.Pointer[slog.Logger]

func SetDefault(logger *slog.Logger) {
	defaultLoggerP.Store(logger)
}
