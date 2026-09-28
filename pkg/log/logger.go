package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
}

type defaultLogger struct {
	*slog.Logger
}

var (
	once sync.Once
	log  Logger = &defaultLogger{}
)

func GetDefaultLogger() Logger {
	once.Do(func() {
		errorOptions := &slog.HandlerOptions{
			Level: slog.LevelError,
		}
		log = &defaultLogger{
			Logger: slog.New(slog.NewTextHandler(os.Stdout, errorOptions)),
		}
	})
	return log
}

func (l *defaultLogger) Infof(format string, args ...any) {
	l.logf(slog.LevelInfo, format, args...)
}

func (l *defaultLogger) Warnf(format string, args ...any) {
	l.logf(slog.LevelWarn, format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	l.logf(slog.LevelError, format, args...)
}

func (l *defaultLogger) Debugf(format string, args ...any) {
	l.logf(slog.LevelDebug, format, args...)
}

// logf formats the message printf-style, as the Logger methods promise
// (slog's own methods would take args as key/value attributes), skipping the
// formatting when level is disabled.
func (l *defaultLogger) logf(level slog.Level, format string, args ...any) {
	ctx := context.Background()
	if !l.Enabled(ctx, level) {
		return
	}
	l.Log(ctx, level, fmt.Sprintf(format, args...))
}
