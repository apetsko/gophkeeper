package grpcserver

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/internal/grpcserver/handlers"
	pb "github.com/apetsko/gophkeeper/protogen/api/proto/v1"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
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

func (s *GRPCServer) BinaryData(ctx context.Context, in *pbrpc.BinaryDataRequest) (*pbrpc.BinaryDataResponse, error) {
	return s.ServerAdmin.BinaryData(ctx, in)
}

func (s *GRPCServer) Records(ctx context.Context, in *pbrpc.RecordsRequest) (*pbrpc.RecordsResponse, error) {
	return s.ServerAdmin.Records(ctx, in)
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
