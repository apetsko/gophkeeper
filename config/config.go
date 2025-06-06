// Package config provides configuration loading and validation for the GophKeeper application.
//
// It supports loading configuration from environment variables or a YAML file, and validates required fields.
// The Config struct holds all application settings, including database, server, JWT, S3, and TLS options.
// The New function loads and validates the configuration, decoding the server encryption key as needed.
package config

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/utils"
	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

// Config holds all application configuration settings, loaded from environment variables or a YAML file.
type Config struct {
	// ConfigFile is the path to the YAML configuration file.
	ConfigFile string `env:"CONFIG_FILE" yaml:"CONFIG_FILE"`
	// DatabaseDSN is the Data Source Name for connecting to the Postgres database.
	DatabaseDSN string `env:"DATABASE_DSN" yaml:"DATABASE_DSN" validate:"required"`
	// GRPCAddress is the address for the gRPC server to listen on.
	GRPCAddress string `env:"GRPC_ADDRESS" yaml:"GRPC_ADDRESS" validate:"required"`
	// HTTPAddress is the address for the HTTP server to listen on.
	HTTPAddress string `env:"HTTP_ADDRESS" yaml:"HTTP_ADDRESS" validate:"required"`
	// StrServerEK is the hex-encoded server encryption key.
	StrServerEK string `env:"SERVER_ENCRYPTION_KEY" yaml:"SERVER_ENCRYPTION_KEY" validate:"required"`
	// ServerEK is the decoded server encryption key as a byte slice.
	ServerEK []byte
	// JWT holds JWT-related configuration.
	JWT JWTConfig `yaml:"JWT"`
	// S3Config holds S3/MinIO-related configuration.
	S3Config S3Config `yaml:"S3"`
	// TLSConfig holds TLS/HTTPS-related configuration.
	TLSConfig TLSConfig `yaml:"TLS"`
}

// JWTConfig contains settings for JWT authentication.
type JWTConfig struct {
	// Secret is the secret key for signing JWT tokens.
	Secret string `env:"JWT_SECRET" yaml:"JWT_SECRET" validate:"required"`
}

// S3Config contains settings for S3/MinIO object storage.
type S3Config struct {
	// AccessKey is the S3 access key.
	AccessKey string `env:"S3_ACCESS_KEY" yaml:"S3_ACCESS_KEY" validate:"required"`
	// SecretKey is the S3 secret key.
	SecretKey string `env:"S3_SECRET_KEY" yaml:"S3_SECRET_KEY" validate:"required"`
	// Bucket is the S3 bucket name.
	Bucket string `env:"S3_BUCKET" yaml:"S3_BUCKET" validate:"required"`
	// Endpoint is the S3 service endpoint.
	Endpoint string `env:"S3_ENDPOINT" yaml:"S3_ENDPOINT" validate:"required"`
}

// TLSConfig contains settings for enabling HTTPS/TLS.
type TLSConfig struct {
	// EnableHTTPS enables HTTPS if set to true.
	EnableHTTPS bool `env:"TLS_ENABLE_HTTPS" yaml:"TLS_ENABLE_HTTPS"`
	// CertPath is the path to the TLS certificate file.
	CertPath string `env:"TLS_CERT_PATH" yaml:"TLS_CERT_PATH"`
	// KeyPath is the path to the TLS private key file.
	KeyPath string `env:"TLS_KEY_PATH" yaml:"TLS_KEY_PATH"`
}

// New loads and validates the application configuration from a YAML file or environment variables.
// It decodes the server encryption key and returns a Config instance.
func New() (*Config, error) {
	var cfg Config

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&cfg.ConfigFile, "f", "", "config.yaml or config.yml")
	_ = fs.Parse(os.Args[1:]) // Не паникуем на ошибке парсинга

	// file or env
	if cfg.ConfigFile != "" {
		slog.Info("Using config file: " + cfg.ConfigFile)
		if err := cfg.readConfigFile(); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	} else {
		if err := env.Parse(&cfg); err != nil {
			return nil, fmt.Errorf("failed to load environment: %w", err)
		}
	}

	// Validate the loaded configuration
	if err := utils.ValidateStruct(cfg); err != nil {
		return nil, err
	}

	serverKey, errDecode := hex.DecodeString(cfg.StrServerEK)
	if errDecode != nil {
		return nil, errDecode
	}

	if len(serverKey) != constants.KeyLength {
		return nil, errors.New("server encryption key must be 32 bytes")
	}

	cfg.ServerEK = serverKey

	return &cfg, nil
}

// readConfigFile loads configuration from the specified YAML file into the Config struct.
func (cfg *Config) readConfigFile() error {
	b, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}
	return nil
}
