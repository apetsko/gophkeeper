package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gophkeeper/internal/grpcserver"
	"gophkeeper/internal/grpcserver/handlers"
	pb "gophkeeper/protogen/api/proto/v1"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	listen, err := net.Listen("tcp", "localhost:3002") // TODO: в конфиг
	if err != nil {
		log.Fatal(err)
	}

	server := grpcserver.NewGRPCServer(handlers.NewServer())

	// Добавляем Unary интерцепторы
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcserver.AuthUnaryInterceptor(map[string]bool{}),
		grpcserver.LoggingInterceptor(),
	))
	pb.RegisterGophKeeperServer(grpcServer, server)
	reflection.Register(grpcServer)

	go func() {
		if grpcErr := grpcServer.Serve(listen); grpcErr != nil {
			log.Panic(fmt.Errorf("grpc server can't start: %w", grpcErr))
		}
	}()

	slog.Info("app is ready")

	<-ctx.Done()
	slog.Info("shutdown signal received")

	grpcServer.GracefulStop()
}
