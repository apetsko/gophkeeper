package config_test

import (
	"testing"

	"github.com/apetsko/gophkeeper/config"
	"github.com/stretchr/testify/require"
)

func TestNew_FromEnv_Success(t *testing.T) {
	t.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")
	t.Setenv("GRPC_ADDRESS", ":50051")
	t.Setenv("HTTP_ADDRESS", ":8080")
	t.Setenv("SERVER_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("S3_ACCESS_KEY", "access")
	t.Setenv("S3_SECRET_KEY", "secret")
	t.Setenv("S3_BUCKET", "bucket")
	t.Setenv("S3_ENDPOINT", "localhost:9000")

	cfg, err := config.New()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 32, len(cfg.ServerEK))
}

func TestNew_InvalidHexKey(t *testing.T) {
	t.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")
	t.Setenv("GRPC_ADDRESS", ":50051")
	t.Setenv("HTTP_ADDRESS", ":8080")
	t.Setenv("SERVER_ENCRYPTION_KEY", "not-hex!")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("S3_ACCESS_KEY", "access")
	t.Setenv("S3_SECRET_KEY", "secret")
	t.Setenv("S3_BUCKET", "bucket")
	t.Setenv("S3_ENDPOINT", "localhost:9000")

	_, err := config.New()
	require.Error(t, err)
}
