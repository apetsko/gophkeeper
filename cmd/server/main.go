package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"gophkeeper/internal/grpcserver"
	"gophkeeper/internal/grpcserver/handlers"
	pb "gophkeeper/protogen/api/proto/v1"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	grpcAddr := ":3007"
	httpAddr := ":8082"

	// GRPC-сервер
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcserver.AuthUnaryInterceptor(map[string]bool{}),
			grpcserver.LoggingInterceptor(),
		),
	)
	pb.RegisterGophKeeperServer(grpcServer, grpcserver.NewGRPCServer(handlers.NewServer()))
	reflection.Register(grpcServer)

	// GRPC Listener
	lis, errListen := net.Listen("tcp", grpcAddr)
	if errListen != nil {
		log.Fatalf("failed to listen: %v", errListen)
	}

	// GRPC запуск
	go func() {
		log.Println("Starting gRPC server on", grpcAddr)
		if errStarting := grpcServer.Serve(lis); errStarting != nil {
			log.Fatalf("gRPC server failed: %v", errStarting)
		}
	}()

	// HTTP Gateway mux
	mux := runtime.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if errRegister := pb.RegisterGophKeeperHandlerFromEndpoint(ctx, mux, grpcAddr, opts); errRegister != nil {
		log.Fatalf("failed to register gRPC-Gateway: %v", errRegister)
	}

	// HTTP сервер
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: handler,
	}

	// HTTP запуск
	go func() {
		log.Println("Starting HTTP server on", httpAddr)
		if errHTTP := httpServer.ListenAndServe(); errHTTP != nil && !errors.Is(errHTTP, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", errHTTP)
		}
	}()

	log.Println("Servers are running...")

	// Ожидаем сигнал завершения
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Контекст с таймаутом для остановки HTTP-сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Остановка gRPC и HTTP серверов
	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	log.Println("Servers stopped gracefully")
}
