package utilities

import (
	"log/slog"
	"os"
	"sync"
)

type Logger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
	Warn(message string, args ...any)
}

var (
	_globalLogger Logger
	_loggerOnce   sync.Once
)

func GetLogger() Logger {
	_loggerOnce.Do(func() {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		_globalLogger = &_DefaultLogger{logger: logger}
	})
	return _globalLogger
}

type _DefaultLogger struct {
	logger *slog.Logger
}

func (l *_DefaultLogger) Info(message string, args ...any) {
	l.logger.Info(message, args...)
}

func (l *_DefaultLogger) Error(message string, args ...any) {
	l.logger.Error(message, args...)
}

func (l *_DefaultLogger) Warn(message string, args ...any) {
	l.logger.Warn(message, args...)
}
