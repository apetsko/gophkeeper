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
	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
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

func (s *GRPCServer) Login(ctx context.Context, in *pbrpcu.LoginRequest) (*pbrpcu.LoginResponse, error) {
	return s.ServerAdmin.Login(ctx, in)
}

func (s *GRPCServer) Signup(ctx context.Context, in *pbrpcu.SignupRequest) (*pbrpcu.SignupResponse, error) {
	return s.ServerAdmin.Signup(ctx, in)
}

func (s *GRPCServer) DataList(ctx context.Context, in *pbrpc.DataListRequest) (*pbrpc.DataListResponse, error) {
	return s.ServerAdmin.DataList(ctx, in)
}

func (s *GRPCServer) DataSave(ctx context.Context, in *pbrpc.DataSaveRequest) (*pbrpc.DataSaveResponse, error) {
	return s.ServerAdmin.DataSave(ctx, in)
}

func (s *GRPCServer) DataDelete(ctx context.Context, in *pbrpc.DataDeleteRequest) (*pbrpc.DataDeleteResponse, error) {
	return s.ServerAdmin.DataDelete(ctx, in)
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
