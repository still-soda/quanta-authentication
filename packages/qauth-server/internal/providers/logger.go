package providers

import (
	"log/slog"
	"os"
)

type ILogger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
	Warn(message string, args ...any)
	With(args ...any) ILogger
}

type RootLogger struct {
	logger *slog.Logger
}

// NewRootLogger 创建一个新的 RootLogger 实例。
func NewRootLogger() *RootLogger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return &RootLogger{logger: logger}
}

func (l *RootLogger) Info(message string, args ...any) {
	l.logger.Info(message, args...)
}

func (l *RootLogger) Error(message string, args ...any) {
	l.logger.Error(message, args...)
}

func (l *RootLogger) Warn(message string, args ...any) {
	l.logger.Warn(message, args...)
}

func (l *RootLogger) With(args ...any) ILogger {
	return &RootLogger{logger: l.logger.With(args...)}
}
