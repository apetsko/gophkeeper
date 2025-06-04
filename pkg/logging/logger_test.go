package logging

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func TestLogger_Methods(t *testing.T) {
	log := NewLogger(slog.LevelDebug)

	testCases := []struct {
		logFunc func()
		name    string
	}{
		{name: "Debug", logFunc: func() { log.Debug("debug message", "key", "value") }},
		{name: "Info", logFunc: func() { log.Info("info message", "key", "value") }},
		{name: "Error", logFunc: func() { log.Error("error message", "key", "value") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logFunc()
		})
	}
}

func TestLogger_FormattedMethods(t *testing.T) {
	log := NewLogger(slog.LevelDebug)
	log.Debugf("debug %s", "message")
	log.Infof("info %s", "message")
	log.Warnf("warn %s", "message")
	log.Errorf("error %s", "message")
}

func TestLogger_FatalMethods(t *testing.T) {
	log := NewLogger(slog.LevelDebug)
	// To avoid os.Exit, run in subprocess
	if os.Getenv("TEST_FATAL") == "1" {
		log.Fatal("fatal message")
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogger_FatalMethods")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); !ok || e.ExitCode() == 0 {
		t.Fatalf("expected exit with non-zero code")
	}
}

func TestLogger_FatalfMethods(t *testing.T) {
	log := NewLogger(slog.LevelDebug)
	if os.Getenv("TEST_FATALF") == "1" {
		log.Fatalf("fatalf %s", "message")
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogger_FatalfMethods")
	cmd.Env = append(os.Environ(), "TEST_FATALF=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); !ok || e.ExitCode() == 0 {
		t.Fatalf("expected exit with non-zero code")
	}
}

func TestInterceptorLogger_AllLevels(t *testing.T) {
	log := NewLogger(slog.LevelDebug)
	logger := InterceptorLogger(log)
	logger.Log(context.Background(), grpc_logging.LevelDebug, "debug", "k", "v")
	logger.Log(context.Background(), grpc_logging.LevelInfo, "info", "k", "v")
	logger.Log(context.Background(), grpc_logging.LevelInfo, "info", "k", "v")
	logger.Log(context.Background(), grpc_logging.LevelError, "error", "k", "v")
}
