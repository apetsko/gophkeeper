package handlers

import (
	"github.com/minio/minio-go/v7"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/crypto"
	"github.com/apetsko/gophkeeper/internal/storage"
)

type ServerAdmin struct {
	storage     *storage.Storage
	jwtConfig   config.JWTConfig
	envelop     *crypto.Envelop
	keyManager  *crypto.KeyManager
	minioBucket string
	minioClient *minio.Client
}

func NewServer(
	storage *storage.Storage,
	jwtConfig config.JWTConfig,
	envelop *crypto.Envelop,
	keyManager *crypto.KeyManager,
	minioBucket string,
	minioClient *minio.Client,
) *ServerAdmin {
	return &ServerAdmin{
		storage:     storage,
		jwtConfig:   jwtConfig,
		envelop:     envelop,
		keyManager:  keyManager,
		minioBucket: minioBucket,
		minioClient: minioClient,
	}
}
