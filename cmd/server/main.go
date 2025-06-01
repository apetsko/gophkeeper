package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/crypto"
	grpcsrv "github.com/apetsko/gophkeeper/internal/server/grpc"
	"github.com/apetsko/gophkeeper/internal/server/grpc/handlers"
	httpsrv "github.com/apetsko/gophkeeper/internal/server/http"
	"github.com/apetsko/gophkeeper/internal/storage"
	"github.com/apetsko/gophkeeper/pkg/logging"
)

func main() {
	slog.SetDefault(logging.LogHandler(slog.LevelDebug))
	log := logging.NewLogger(slog.LevelDebug)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config read err %v", err)
	}

	dbClient, err := storage.NewPostrgesClient(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database client init err %v", err)
	}

	minioClient, err := storage.NewMinioClient(ctx, cfg.Minio)
	if err != nil {
		log.Fatalf("minio client init err %v", err)
	}

	envelop := crypto.NewEnvelop(dbClient)
	keyManager := crypto.NewKeyManager(dbClient, cfg.ServerEK)

	sa := handlers.ServerAdmin{
		Storage:     dbClient,
		JWTConfig:   cfg.JWT,
		Envelop:     envelop,
		KeyManager:  keyManager,
		MinioBucket: cfg.Minio.Bucket,
		MinioClient: minioClient,
	}

	// Start gRPC server
	if _, err := grpcsrv.RunGRPC(cfg, &sa, log); err != nil {
		log.Fatal("gRPC server failed: " + err.Error())
	}

	// Start HTTP server
	if _, err := httpsrv.RunHTTP(ctx, cfg, log); err != nil {
		log.Fatal("HTTP server failed: " + err.Error())
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	log.Info("Shutting down servers...")
}
