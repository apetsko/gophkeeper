package config

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/apetsko/gophkeeper/utils"
	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigFile  string `env:"CONFIG_FILE" yaml:"CONFIG_FILE"`
	DatabaseDSN string `env:"DATABASE_DSN" yaml:"DATABASE_DSN" validate:"required"`
	GRPCAddress string `env:"GRPC_ADDRESS" yaml:"GRPC_ADDRESS" validate:"required"`
	HTTPAddress string `env:"HTTP_ADDRESS" yaml:"HTTP_ADDRESS" validate:"required"`
	StrServerEK string `env:"SERVER_ENCRYPTION_KEY" yaml:"SERVER_ENCRYPTION_KEY" validate:"required"`
	ServerEK    []byte
	JWT         JWTConfig   `yaml:"JWT"`
	Minio       MinioConfig `yaml:"MINIO"`
	TLSConfig   TLSConfig   `yaml:"TLS"`
}
type JWTConfig struct {
	Secret string `env:"JWT_SECRET" yaml:"JWT_SECRET" validate:"required"`
}
type MinioConfig struct {
	ID      string `env:"MINIO_ID" yaml:"MINIO_ID" validate:"required"`
	Secret  string `env:"MINIO_SECRET" yaml:"MINIO_SECRET" validate:"required"`
	Bucket  string `env:"MINIO_BUCKET" yaml:"MINIO_BUCKET" validate:"required"`
	Address string `env:"MINIO_ADDRESS" yaml:"MINIO_ADDRESS" validate:"required"`
}

type TLSConfig struct {
	EnableHTTPS bool   `env:"TLS_ENABLE_HTTPS" yaml:"TLS_ENABLE_HTTPS"`
	CertPath    string `env:"TLS_CERT_PATH" yaml:"TLS_CERT_PATH"`
	KeyPath     string `env:"TLS_KEY_PATH" yaml:"TLS_KEY_PATH"`
}

func New() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.ConfigFile, "f", "", "config.yaml or config.yml")
	flag.Parse()

	//file or env
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

	if len(serverKey) != 32 {
		return nil, errors.New("server encryption key must be 32 bytes")
	}

	cfg.ServerEK = serverKey

	return &cfg, nil
}

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
