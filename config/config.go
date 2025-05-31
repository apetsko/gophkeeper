package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/apetsko/gophkeeper/utils"
	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigFile         string      `env:"CONFIG_FILE" yaml:"CONFIG_FILE"`
	DatabaseDSN        string      `env:"DATABASE_DSN" yaml:"DATABASE_DSN" validate:"required"`
	GRPCAddress        string      `env:"GRPC_ADDRESS" yaml:"GRPC_ADDRESS" validate:"required"`
	GRPCGatewayAddress string      `env:"GRPC_GATEWAY_ADDRESS" yaml:"GRPC_GATEWAY_ADDRESS" validate:"required"`
	JWT                JWTConfig   `yaml:"JWT"`
	Minio              MinioConfig `yaml:"MINIO"`
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
