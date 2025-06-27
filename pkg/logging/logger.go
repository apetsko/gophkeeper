// Package logging provides structured logging utilities and wrappers around slog,
// including convenience methods and gRPC middleware integration.
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

// GenerateJWT creates a signed JWT token for the given user ID and username.
//
// The token uses HS256 signing and includes user ID, username, and issued-at claims.
//
// Returns the signed JWT string or an error if signing fails.
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

// Logger wraps slog.Logger and provides convenience methods for formatted logging.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new Logger with the specified log level.
func NewLogger(level slog.Level) *Logger {
	base := LogHandler(level)
	return &Logger{Logger: base}
}

// Debugf logs a debug-level message with formatting.
func (l *Logger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof logs an info-level message with formatting.
func (l *Logger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warnf logs a warning-level message with formatting.
func (l *Logger) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Errorf logs an error-level message with formatting.
func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

// Fatal logs an error-level message and exits the application.
func (l *Logger) Fatal(msg string) {
	l.Log(context.Background(), slog.LevelError, msg)
	os.Exit(1)
}

// Fatalf logs a formatted error-level message and exits the application.
func (l *Logger) Fatalf(format string, args ...any) {
	l.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// InterceptorLogger returns a grpc_logging.Logger compatible with slog for gRPC middleware.
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
