package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

	// Канал для синхронизации запуска gRPC-сервера
	grpcReady := make(chan struct{})

	// Запуск gRPC-сервера в отдельной горутине
	go func() {
		listen, err := net.Listen("tcp", ":3007")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		server := grpcserver.NewGRPCServer(handlers.NewServer())

		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				grpcserver.AuthUnaryInterceptor(map[string]bool{}),
				grpcserver.LoggingInterceptor(),
			),
		)

		pb.RegisterGophKeeperServer(grpcServer, server)
		reflection.Register(grpcServer)

		// Сигнализируем, что gRPC-сервер готов
		close(grpcReady)

		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	// Ожидаем, пока gRPC-сервер будет готов
	<-grpcReady

	// Инициализация gRPC-Gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterGophKeeperHandlerFromEndpoint(ctx, mux, ":3007", opts); err != nil {
		log.Fatalf("failed to register gRPC-Gateway: %v", err)
	}

	go func() {
		if err := http.ListenAndServe(":8082", mux); err != nil {
			log.Fatalf("failed to serve HTTP server: %v", err)
		}
	}()

	log.Println("Servers are running...")

	<-ctx.Done()
	log.Println("Shutdown signal received")
}
