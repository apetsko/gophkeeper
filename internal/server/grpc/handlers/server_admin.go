package handlers

import (
	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/crypto"
	"github.com/apetsko/gophkeeper/internal/storage"
)

type ServerAdmin struct {
	Storage    storage.IStorage
	StorageS3  storage.S3Client
	JWTConfig  config.JWTConfig
	Envelope   crypto.IEnvelope
	KeyManager crypto.KeyManagerInterface
}

func NewServerAdmin(
	storage storage.IStorage,
	storageS3 storage.S3Client,
	jwtConfig config.JWTConfig,
	envelope crypto.IEnvelope,
	keyManager crypto.KeyManagerInterface,
) *ServerAdmin {
	return &ServerAdmin{
		Storage:    storage,
		StorageS3:  storageS3,
		JWTConfig:  jwtConfig,
		Envelope:   envelope,
		KeyManager: keyManager,
	}
}
