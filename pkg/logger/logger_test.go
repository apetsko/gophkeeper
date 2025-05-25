package logger

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogger_Methods(t *testing.T) {
	logger, _ := New(zapcore.DebugLevel)

	testCases := []struct {
		logFunc func()
		name    string
	}{
		{name: "Debug", logFunc: func() { logger.Debug("debug message", "key", "value") }},
		{name: "Info", logFunc: func() { logger.Info("info message", "key", "value") }},
		{name: "Error", logFunc: func() { logger.Error("error message", "key", "value") }},
		{name: "Printf", logFunc: func() { logger.Printf("formatted %s", "message") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logFunc()
		})
	}
}
