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
		log.Errorf("config read err %v", err)
		return
	}

	dbClient, err := storage.NewPostgresClient(cfg.DatabaseDSN)
	if err != nil {
		log.Errorf("database client init err %v", err)
		return
	}

	s3Client, err := storage.NewS3Client(ctx, cfg.S3Config)
	if err != nil {
		log.Errorf("minio client init err %v", err)
		return
	}

	envelope := crypto.NewEnvelope(dbClient)
	keyManager := crypto.NewKeyManager(dbClient, cfg.ServerEK)

	sa := handlers.NewServerAdmin(dbClient, s3Client, cfg.JWT, envelope, keyManager)

	// Start gRPC server
	if _, err := grpcsrv.RunGRPC(cfg, sa, log); err != nil {
		log.Errorf("gRPC server failed: %v", err.Error())
		return
	}

	// Start HTTP server
	if _, err := httpsrv.RunHTTP(ctx, cfg, log); err != nil {
		log.Errorf("HTTP server failed:%v ", err.Error())
		return
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	log.Info("Shutting down servers...")
}
