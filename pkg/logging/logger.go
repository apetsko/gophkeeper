// Package logging provides structured logging functionality for the application.
// It wraps slog to enable easy and efficient logging with support
// for different log levels and structured log entries.
package logging

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// LogEntry defines the interface for log entries.
type LogEntry interface {
	// Write logs the status, bytes, header, elapsed time, and extra information.
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{})
	// Panic logs a panic message with the stack trace.
	Panic(v interface{}, stack []byte)
}

// Logger wraps slog.Logger to provide structured logging.
type Logger struct {
	slog *slog.Logger
}

// New создает новый Logger с JSON-форматом, временем в RFC3339 с миллисекундами,
// и без поля source (caller).
func New(level slog.Level) *Logger {
	//_, file, line, ok := runtime.Caller(2) // Skip 2 frames to get the caller of logWithLine
	//if !ok {
	//	file = "unknown"
	//	line = 0
	//}
	//var allArgs []any
	//allArgs = append(allArgs, "file", file, "line", line)
	//allArgs = append(allArgs, args...)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false, // Отключаем вывод source (caller)
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					// Формат с миллисекундами RFC3339Nano (потому что RFC3339 неявно не выводит миллисекунды, RFC3339Nano - с наносекундами)
					// Но можно ограничить наносекунды до миллисекунд через Format с 000:
					return slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(t.Format("2006-01-02T15:04:05.000Z07:00")),
					}
				}
			}
			return a
		},
	})

	return &Logger{slog: slog.New(handler)}
}

// Close is a no-op for slog (included for compatibility).
func (l *Logger) Close() error {
	// slog doesn't need explicit sync like zap.
	return nil
}

// Debug logs a debug message with additional context.
func (l *Logger) Debug(message string, keysAndValues ...interface{}) {
	l.slog.Debug(message, keysAndValues...)
}

// Info logs an informational message with additional context.
func (l *Logger) Info(message string, keysAndValues ...interface{}) {
	l.slog.Info(message, keysAndValues...)
}

// Error logs an error message with additional context.
func (l *Logger) Error(message string, keysAndValues ...interface{}) {
	l.slog.Error(message, keysAndValues...)
}

// Fatal logs a fatal message and exits the application.
func (l *Logger) Fatal(message string, keysAndValues ...interface{}) {
	l.slog.Error(message, keysAndValues...)
	os.Exit(1)
}

// Printf logs a formatted informational message.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.slog.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.slog.Error(fmt.Sprintf(format, args...))
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
			l.slog.Warn(msg, args...)
		case grpc_logging.LevelError:
			l.Error(msg, args...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
