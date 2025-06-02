package handlers

import (
	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/crypto"
	"github.com/apetsko/gophkeeper/internal/storage"
)

type ServerAdmin struct {
	Storage    *storage.Storage
	StorageS3  *storage.S3
	JWTConfig  config.JWTConfig
	Envelop    *crypto.Envelop
	KeyManager *crypto.KeyManager
}
