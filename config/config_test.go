package config_test

import (
	"os"
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

func TestNew_FromFile_Success(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := `
DATABASE_DSN: "postgres://user:pass@localhost/db"
GRPC_ADDRESS: ":50051"
HTTP_ADDRESS: ":8080"
SERVER_ENCRYPTION_KEY: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
JWT:
  JWT_SECRET: "secret"
S3:
  S3_ACCESS_KEY: "access"
  S3_SECRET_KEY: "secret"
  S3_BUCKET: "bucket"
  S3_ENDPOINT: "localhost:9000"
`
	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	// Simulate passing -f flag
	os.Args = []string{"test", "-f", tmpFile.Name()}

	cfg, err := config.New()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 32, len(cfg.ServerEK))
}

func TestNew_ServerKeyWrongLength(t *testing.T) {
	t.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")
	t.Setenv("GRPC_ADDRESS", ":50051")
	t.Setenv("HTTP_ADDRESS", ":8080")
	// 16 bytes instead of 32
	t.Setenv("SERVER_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("S3_ACCESS_KEY", "access")
	t.Setenv("S3_SECRET_KEY", "secret")
	t.Setenv("S3_BUCKET", "bucket")
	t.Setenv("S3_ENDPOINT", "localhost:9000")

	_, err := config.New()
	require.ErrorContains(t, err, "server encryption key must be 32 bytes")
}

func TestNew_ConfigFileNotFound(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"test", "-f", "/nonexistent/path.yaml"}
	_, err := config.New()
	require.ErrorContains(t, err, "failed to load config from file")
}

func TestNew_MissingRequiredEnv(t *testing.T) {
	// Do not set DATABASE_DSN
	t.Setenv("GRPC_ADDRESS", ":50051")
	t.Setenv("HTTP_ADDRESS", ":8080")
	t.Setenv("SERVER_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("S3_ACCESS_KEY", "access")
	t.Setenv("S3_SECRET_KEY", "secret")
	t.Setenv("S3_BUCKET", "bucket")
	t.Setenv("S3_ENDPOINT", "localhost:9000")

	_, err := config.New()
	require.ErrorContains(t, err, "failed on the 'required' tag")
}
