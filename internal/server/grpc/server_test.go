package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/internal/server/grpc/handlers"
	pb "github.com/apetsko/gophkeeper/protogen/api/proto/v1"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func dialer(srv *grpc.Server) func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(bufSize)
	go func() {
		_ = srv.Serve(lis)
	}()
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestGRPCHandler_Ping(t *testing.T) {
	// Mocks and config
	st := mocks.NewIStorage(t)
	s3 := mocks.NewS3Client(t)
	env := mocks.NewIEnvelope(t)
	km := mocks.NewKeyManagerInterface(t)
	cfg := config.JWTConfig{Secret: "testsecret"}

	admin := handlers.NewServerAdmin(st, s3, cfg, env, km)
	srv := grpc.NewServer()
	pb.RegisterGophKeeperServer(srv, NewGRPCHandler(admin))

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(srv)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGophKeeperClient(conn)
	resp, err := client.Ping(ctx, &pbrpc.PingRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
}
