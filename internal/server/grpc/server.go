// Package grpc provides the gRPC server implementation for the GophKeeper service.
//
// This package sets up the gRPC server, configures authentication middleware,
// registers service handlers, and manages server lifecycle and TLS settings.
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

// GRPCHandler implements the gRPC server interface for the GophKeeper service.
//
// This struct embeds ServerAdmin to delegate business logic for user and data operations.
// It provides gRPC method handlers for health checks, user authentication, registration,
// and CRUD operations on user data records.
type GRPCHandler struct {
	pb.UnimplementedGophKeeperServer
	*handlers.ServerAdmin
}

// NewGRPCHandler creates a new GRPCHandler instance.
//
// Parameters:
//   - admin: Pointer to ServerAdmin containing business logic.
//
// Returns:
//   - pb.GophKeeperServer: The gRPC server handler.
func NewGRPCHandler(admin *handlers.ServerAdmin) pb.GophKeeperServer {
	return &GRPCHandler{
		ServerAdmin: admin,
	}
}

// Ping handles a health check request for the gRPC service.
//
// This method verifies that the server is alive and responding to requests.
// It delegates the call to ServerAdmin and always returns a successful response.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The PingRequest message.
//
// Returns:
//   - *pbrpc.PingResponse: An empty response indicating success.
//   - error: Always nil.
func (s *GRPCHandler) Ping(ctx context.Context, in *pbrpc.PingRequest) (*pbrpc.PingResponse, error) {
	return s.ServerAdmin.Ping(ctx, in)
}

// Login handles the gRPC request for user authentication.
//
// This method validates the username and password, checks credentials against the database,
// generates a JWT token upon successful authentication, and ensures the user's master key exists.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The LoginRequest message with user credentials.
//
// Returns:
//   - *pbrpcu.LoginResponse: User details and authentication token.
//   - error: An error if authentication fails.
func (s *GRPCHandler) Login(ctx context.Context, in *pbrpcu.LoginRequest) (*pbrpcu.LoginResponse, error) {
	return s.ServerAdmin.Login(ctx, in)
}

func (s *GRPCHandler) Signup(ctx context.Context, in *pbrpcu.SignupRequest) (*pbrpcu.SignupResponse, error) {
	return s.ServerAdmin.Signup(ctx, in)
}

// DataList handles the gRPC request to list all user data records.
//
// This method checks user authorization and retrieves a list of data records
// associated with the authenticated user.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The DataListRequest message.
//
// Returns:
//   - *pbrpc.DataListResponse: List of user data records.
//   - error: A gRPC error if access is denied or an internal error occurs.
func (s *GRPCHandler) DataList(ctx context.Context, in *pbrpc.DataListRequest) (*pbrpc.DataListResponse, error) {
	return s.ServerAdmin.DataList(ctx, in)
}

// DataSave handles the gRPC request to create or update a user data record.
//
// This method checks user authorization, validates the input, encrypts the data,
// and saves it to the database or storage.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The DataSaveRequest message with data to save.
//
// Returns:
//   - *pbrpc.DataSaveResponse: Confirmation of save operation.
//   - error: A gRPC error if access is denied or an internal error occurs.
func (s *GRPCHandler) DataSave(ctx context.Context, in *pbrpc.DataSaveRequest) (*pbrpc.DataSaveResponse, error) {
	return s.ServerAdmin.DataSave(ctx, in)
}

// DataDelete handles the gRPC request to delete a user's data record.
//
// This method checks user authorization, verifies ownership of the data record,
// and deletes the record from storage if permitted.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The DataDeleteRequest message with the data record ID.
//
// Returns:
//   - *pbrpc.DataDeleteResponse: Success message if deletion is successful.
//   - error: A gRPC error if the user is not authorized or an internal error occurs.
func (s *GRPCHandler) DataDelete(ctx context.Context, in *pbrpc.DataDeleteRequest) (*pbrpc.DataDeleteResponse, error) {
	return s.ServerAdmin.DataDelete(ctx, in)
}

// DataView handles the gRPC request to retrieve a specific user data record by its ID.
//
// This method checks user authorization, fetches the encrypted data from storage,
// decrypts it, and returns the result in the response.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The DataViewRequest message containing the record ID.
//
// Returns:
//   - *pbrpc.DataViewResponse: The requested user data record.
//   - error: A gRPC error if access is denied or an internal error occurs.
func (s *GRPCHandler) DataView(ctx context.Context, in *pbrpc.DataViewRequest) (*pbrpc.DataViewResponse, error) {
	return s.ServerAdmin.DataView(ctx, in)
}

// RunGRPC starts the gRPC server for the GophKeeper service.
//
// This function configures the gRPC server with optional TLS, authentication middleware,
// logging interceptors, and registers the service handlers. It listens on the configured
// address and manages the server lifecycle using an error group.
//
// Parameters:
//   - cfg: Pointer to the application configuration.
//   - sa: Pointer to ServerAdmin containing business logic.
//   - log: Logger instance for server logging.
//
// Returns:
//   - *grpc.Server: The running gRPC server instance.
//   - error: An error if the server fails to start.
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
				"/api.proto.v1.GophKeeper/DataView":   true,
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
			log.Error("gRPC server error: " + err.Error())
		}
	}()

	return srv, nil
}

// authUnaryInterceptor returns a gRPC unary server interceptor for JWT authentication.
//
// This interceptor checks if the called method requires authentication. For protected methods,
// it extracts and validates the JWT from the request metadata, verifies the signing method and claims,
// and injects the user ID and JWT into the context for downstream handlers.
//
// Parameters:
//   - protected: Map of gRPC method names that require authentication.
//   - jwtSecret: Secret key used to validate JWT tokens.
//
// Returns:
//   - grpc.UnaryServerInterceptor: The configured authentication interceptor.
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
