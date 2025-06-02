package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/golang-jwt/jwt/v5"
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/internal/server/grpc/handlers"
	"github.com/apetsko/gophkeeper/pkg/logging"
	pb "github.com/apetsko/gophkeeper/protogen/api/proto/v1"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
)

type GRPCHandler struct {
	pb.UnimplementedGophKeeperServer
	*handlers.ServerAdmin
}

func NewGRPCHandler(admin *handlers.ServerAdmin) pb.GophKeeperServer {
	return &GRPCHandler{
		ServerAdmin: admin,
	}
}

func (s *GRPCHandler) Ping(ctx context.Context, in *pbrpc.PingRequest) (*pbrpc.PingResponse, error) {
	return s.ServerAdmin.Ping(ctx, in)
}

func (s *GRPCHandler) Login(ctx context.Context, in *pbrpcu.LoginRequest) (*pbrpcu.LoginResponse, error) {
	return s.ServerAdmin.Login(ctx, in)
}

func (s *GRPCHandler) Signup(ctx context.Context, in *pbrpcu.SignupRequest) (*pbrpcu.SignupResponse, error) {
	return s.ServerAdmin.Signup(ctx, in)
}

func (s *GRPCHandler) DataList(ctx context.Context, in *pbrpc.DataListRequest) (*pbrpc.DataListResponse, error) {
	return s.ServerAdmin.DataList(ctx, in)
}

func (s *GRPCHandler) DataSave(ctx context.Context, in *pbrpc.DataSaveRequest) (*pbrpc.DataSaveResponse, error) {
	return s.ServerAdmin.DataSave(ctx, in)
}

func (s *GRPCHandler) DataDelete(ctx context.Context, in *pbrpc.DataDeleteRequest) (*pbrpc.DataDeleteResponse, error) {
	return s.ServerAdmin.DataDelete(ctx, in)
}

func (s *GRPCHandler) DataView(ctx context.Context, in *pbrpc.DataViewRequest) (*pbrpc.DataViewResponse, error) {
	return s.ServerAdmin.DataView(ctx, in)
}

func authUnaryInterceptor(protected map[string]bool, jwtSecret []byte) grpc.UnaryServerInterceptor {
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

// RunGRPC запускает gRPC сервер
func RunGRPC(cfg *config.Config, sa *handlers.ServerAdmin, log *logging.Logger) (*grpc.Server, error) {

	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.GRPCAddress, err)
	}

	var opts []grpc.ServerOption

	if cfg.TLSConfig.EnableHTTPS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSConfig.CertPath, cfg.TLSConfig.KeyPath)
		if err != nil {
			log.Fatalf("failed to load TLS credentials: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(
		authUnaryInterceptor(
			map[string]bool{
				"/api.proto.v1.GophKeeper/DataList":   true,
				"/api.proto.v1.GophKeeper/DataSave":   true,
				"/api.proto.v1.GophKeeper/DataDelete": true,
			},
			[]byte(cfg.JWT.Secret),
		),
		grpcLogging.UnaryServerInterceptor(logging.InterceptorLogger(log)),
	))

	srv := grpc.NewServer(opts...)

	h := NewGRPCHandler(sa)

	pb.RegisterGophKeeperServer(srv, h)
	reflection.Register(srv)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		<-ctx.Done()
		srv.GracefulStop()
		return nil
	})

	g.Go(func() error {
		log.Info(fmt.Sprintf("Starting gRPC server at %s, TLS: %t", cfg.GRPCAddress, cfg.TLSConfig.EnableHTTPS))
		return srv.Serve(lis)
	})

	go func() {
		if err := g.Wait(); err != nil {
			log.Error("gRPC server error", err)
		}
	}()

	return srv, nil
}
