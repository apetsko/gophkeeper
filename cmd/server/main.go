// Package main is the entry point for the GophKeeper server application.
//
// This server initializes configuration, logging, database (Postgres), and object storage (MinIO/S3) clients.
// It sets up cryptographic services, JWT authentication, and starts both gRPC and HTTP servers.
// The application supports graceful shutdown on SIGTERM, SIGINT, or SIGQUIT signals.
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
	"github.com/apetsko/gophkeeper/pkg/version"
)

// main initializes all server dependencies and starts the gRPC and HTTP servers.
// It handles configuration loading, logging setup, database and S3 client creation,
// cryptographic envelope and key manager setup, and graceful shutdown.
func main() {
	version.PrintVersion()

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
