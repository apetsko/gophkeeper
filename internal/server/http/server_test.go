package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/internal/server/grpc"
	"github.com/apetsko/gophkeeper/internal/server/grpc/handlers"
	"github.com/apetsko/gophkeeper/pkg/logging"
)

func TestHTTPServer_Ping(t *testing.T) {
	st := mocks.NewIStorage(t)
	s3 := mocks.NewS3Client(t)
	env := mocks.NewIEnvelope(t)
	km := mocks.NewKeyManagerInterface(t)
	jwtCfg := config.JWTConfig{Secret: "testsecret"}
	admin := handlers.NewServerAdmin(st, s3, jwtCfg, env, km)

	cfg := &config.Config{
		GRPCAddress: "127.0.0.1:9090",
		HTTPAddress: "127.0.0.1:8080",
		JWT:         jwtCfg,
		TLSConfig:   config.TLSConfig{EnableHTTPS: false},
	}
	log := logging.NewLogger(slog.LevelDebug)

	grpcSrv, err := grpc.RunGRPC(cfg, admin, log)
	if err != nil {
		t.Fatalf("failed to start gRPC server: %v", err)
	}
	defer grpcSrv.GracefulStop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	httpSrv, err := RunHTTP(ctx, cfg, log)
	if err != nil {
		t.Fatalf("failed to start HTTP server: %v", err)
	}
	defer httpSrv.Close()

	// Wait for the HTTP server to start
	time.Sleep(200 * time.Millisecond)

	// Extract the actual port the HTTP server is listening on
	addr := httpSrv.Addr
	if strings.HasSuffix(addr, ":0") {
		// If Addr is still :0, get the actual port from the listener
		addr = httpSrv.Addr
		if addr == "" {
			t.Fatal("HTTP server Addr is empty")
		}
		// Try to resolve the port from the listener
		lnAddr := httpSrv.Addr
		if lnAddr == "" {
			t.Fatal("could not determine HTTP server address")
		}
		addr = lnAddr
	}

	// If Addr is ":0", get the actual port from the listener
	if strings.HasSuffix(addr, ":0") {
		t.Fatal("HTTP server is still listening on :0")
	}

	url := fmt.Sprintf("http://%s/v1/ping", addr)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to GET /v1/ping: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["message"] != "pong" {
		t.Errorf("unexpected message: %v", result["message"])
	}
}
