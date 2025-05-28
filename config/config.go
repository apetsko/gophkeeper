package config

import (
	"flag"
	"os"
)

type MinioConfig struct {
	ID      string `env:"MINIO_ID"`
	Secret  string `env:"MINIO_SECRET"`
	Bucket  string `env:"MINIO_BUCKET"`
	Address string `env:"MINIO_ADDRESS"`
}

type Config struct {
	DatabaseDSN        string `env:"DATABASE_DSN"`
	GRPCAddress        string `env:"GRPC_ADDRESS"`
	GRPCGatewayAddress string `env:"GRPC_GATEWAY_ADDRESS"`
	Minio              MinioConfig
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres:postgres@localhost:25432/gophkeeper?sslmode=disable", "database DSN")
	flag.StringVar(&cfg.GRPCAddress, "g", ":3007", "GRPC server startup address")
	flag.StringVar(&cfg.GRPCGatewayAddress, "h", ":8082", "GRPCGateway startup address")

	flag.Parse()

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		cfg.DatabaseDSN = databaseDSN
	}

	if gRPCAddress := os.Getenv("GRPC_ADDRESS"); gRPCAddress != "" {
		cfg.GRPCAddress = gRPCAddress
	}

	if gRPCGatewayAddress := os.Getenv("GRPC_GATEWAY_ADDRESS"); gRPCGatewayAddress != "" {
		cfg.GRPCGatewayAddress = gRPCGatewayAddress
	}

	if minioId := os.Getenv("MINIO_ID"); minioId != "" {
		cfg.Minio.ID = minioId
	} else {
		cfg.Minio.ID = "minioadmin"
	}

	if minioSecret := os.Getenv("MINIO_SECRET"); minioSecret != "" {
		cfg.Minio.Secret = minioSecret
	} else {
		cfg.Minio.Secret = "minioadmin"
	}

	if minioBucket := os.Getenv("MINIO_BUCKET"); minioBucket != "" {
		cfg.Minio.Bucket = minioBucket
	} else {
		cfg.Minio.Bucket = "gophkeeper"
	}

	if minioAddress := os.Getenv("MINIO_ADDRESS"); minioAddress != "" {
		cfg.Minio.Address = minioAddress
	} else {
		cfg.Minio.Address = "localhost:9000"
	}

	return &cfg, nil
}
