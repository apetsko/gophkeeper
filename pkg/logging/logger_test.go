package logging

import (
	"log/slog"
	"testing"
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
