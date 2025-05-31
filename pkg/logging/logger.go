// Package logging provides structured logging functionality for the application.
// It wraps slog to enable easy and efficient logging with support
// for different log levels and structured log entries.
package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func LogHandler(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a = slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(t.Format("2006-01-02T15:04:05.000Z07:00")),
					}
				}
			}

			if a.Key == slog.SourceKey {
				src, ok := a.Value.Any().(*slog.Source)
				if ok {
					link := fmt.Sprintf("file://%s:%d", src.File, src.Line)
					a = slog.Attr{
						Key:   filepath.Base(src.Function),
						Value: slog.StringValue(link),
					}
				}
			}
			return a
		},
	})
	return slog.New(handler)
}

type Logger struct {
	*slog.Logger
}

func NewLogger(level slog.Level) *Logger {
	base := LogHandler(level)
	return &Logger{Logger: base}
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(msg string) {
	l.Log(context.Background(), slog.LevelError, msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// InterceptorLogger returns a grpc-compatible slog logger.
func InterceptorLogger(l *Logger) grpc_logging.Logger {
	return grpc_logging.LoggerFunc(func(ctx context.Context, lvl grpc_logging.Level, msg string, fields ...any) {
		args := make([]any, 0, len(fields))

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]
			args = append(args, key, value)
		}

		switch lvl {
		case grpc_logging.LevelDebug:
			l.Debug(msg, args...)
		case grpc_logging.LevelInfo:
			l.Info(msg, args...)
		case grpc_logging.LevelWarn:
			slog.Warn(msg, args...)
		case grpc_logging.LevelError:
			l.Error(msg, args...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
