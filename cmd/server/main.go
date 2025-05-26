package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/apetsko/gophkeeper/internal/grpcserver"
	"github.com/apetsko/gophkeeper/internal/grpcserver/handlers"
	"github.com/apetsko/gophkeeper/pkg/logging"
	pb "github.com/apetsko/gophkeeper/protogen/api/proto/v1"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logging.New(slog.LevelDebug)
	defer log.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	grpcAddr := ":3007"
	httpAddr := ":8082"

	// GRPC-сервер
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcserver.AuthUnaryInterceptor(map[string]bool{}),                  // твой авторизационный interceptor
			grpc_logging.UnaryServerInterceptor(logging.InterceptorLogger(log)), // логгер
		),
	)

	pb.RegisterGophKeeperServer(grpcServer, grpcserver.NewGRPCServer(handlers.NewServer()))
	reflection.Register(grpcServer)

	// Запускаем сервера
	go runGRPC(ctx, grpcServer, grpcAddr, log)
	go runHTTP(ctx, httpAddr, grpcAddr, log)

	log.Info("Servers are running...")

	// Ждём сигнал завершения
	<-ctx.Done()
	log.Info("Shutdown signal received")

	// Контекст с таймаутом для graceful shutdown HTTP
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем gRPC и HTTP серверы
	grpcServer.GracefulStop()
	if err := shutdownHTTP(shutdownCtx, httpAddr, log); err != nil {
		log.Error("HTTP shutdown error", "err", err)
	}

	log.Info("Servers stopped gracefully")
}

// runGRPC запускает gRPC сервер
func runGRPC(ctx context.Context, grpcServer *grpc.Server, grpcAddr string, log *logging.Logger) {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}

	log.Info("Starting gRPC server on " + grpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

// runHTTP запускает HTTP сервер с gRPC-Gateway и CORS
func runHTTP(ctx context.Context, httpAddr, grpcAddr string, log *logging.Logger) {
	mux := runtime.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterGophKeeperHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		log.Fatalf("failed to register gRPC-Gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: handler,
	}

	log.Info("Starting HTTP server on " + httpAddr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

// shutdownHTTP выполняет корректное завершение HTTP сервера
func shutdownHTTP(ctx context.Context, httpAddr string, log *logging.Logger) error {
	httpServer := &http.Server{Addr: httpAddr}
	err := httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
