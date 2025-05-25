package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/grpcserver/handlers"
	pb "gophkeeper/protogen/api/proto/v1"
	pbrpc "gophkeeper/protogen/api/proto/v1/rpc"
)

type GRPCServer struct {
	pb.UnimplementedGophKeeperServer
	*handlers.ServerAdmin
}

func NewGRPCServer(admin *handlers.ServerAdmin) pb.GophKeeperServer {
	return &GRPCServer{
		ServerAdmin: admin,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *pbrpc.PingRequest) (*pbrpc.PingResponse, error) {
	return s.ServerAdmin.Ping(ctx, in)
}

func (s *GRPCServer) Credentials(ctx context.Context, in *pbrpc.CredentialsRequest) (*pbrpc.CredentialsResponse, error) {
	return s.ServerAdmin.Credentials(ctx, in)
}

func (s *GRPCServer) BankCard(ctx context.Context, in *pbrpc.BankCardRequest) (*pbrpc.BankCardResponse, error) {
	return s.ServerAdmin.BankCard(ctx, in)
}

func (s *GRPCServer) BinaryData(stream pb.GophKeeper_BinaryDataServer) error {
	return s.ServerAdmin.BinaryData(stream)
}

func AuthUnaryInterceptor(protected map[string]bool) grpc.UnaryServerInterceptor {
	slog.Info("Auth interceptor enabled")

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !protected[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		jwt := md.Get(string(constants.JWT))
		if len(jwt) == 0 || jwt[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "missing jwt")
		}

		ctx = context.WithValue(ctx, constants.JWT, jwt[0])

		return handler(ctx, req)
	}
}

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	slog.Info("gRPC logger interceptor enabled")

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st := status.Convert(err)

		message := fmt.Sprintf("method: %s, duration: %s, status: %s",
			info.FullMethod,
			duration,
			st.Code().String(),
		)

		if err != nil {
			message += fmt.Sprintf(", error: %s", st.Message())
		}

		slog.Info(message)

		return resp, err
	}
}
