package grpcserver

import (
	"context"
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
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

func AuthUnaryInterceptor(protected map[string]bool, jwtSecret []byte) grpc.UnaryServerInterceptor {
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

		jwtHeader := md.Get(string(constants.JWT))
		if len(jwtHeader) == 0 || jwtHeader[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "missing jwt")
		}

		tokenStr := jwtHeader[0]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid jwt")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if uidFloat, ok := claims["user_id"].(float64); ok {
				userID := int(uidFloat)
				ctx = context.WithValue(ctx, constants.UserID, userID)
			} else {
				return nil, status.Error(codes.InvalidArgument, "user_id not found or not a number")
			}
		}

		ctx = context.WithValue(ctx, constants.JWT, tokenStr)
		return handler(ctx, req)
	}
}
