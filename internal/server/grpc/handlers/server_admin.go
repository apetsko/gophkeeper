package handlers

import (
	"github.com/minio/minio-go/v7"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/crypto"
	"github.com/apetsko/gophkeeper/internal/storage"
)

type ServerAdmin struct {
	Storage     *storage.Storage
	JWTConfig   config.JWTConfig
	Envelop     *crypto.Envelop
	KeyManager  *crypto.KeyManager
	MinioBucket string
	MinioClient *minio.Client
}
