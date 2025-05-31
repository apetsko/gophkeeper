package handlers

import (
	"github.com/minio/minio-go/v7"

	"github.com/apetsko/gophkeeper/internal/storage"
)

type ServerAdmin struct {
	dbClient    *storage.Storage
	minioBucket string
	minioClient *minio.Client
}

func NewServer(
	dbClient *storage.Storage,
	minioBucket string,
	minioClient *minio.Client,
) *ServerAdmin {
	return &ServerAdmin{
		dbClient:    dbClient,
		minioBucket: minioBucket,
		minioClient: minioClient,
	}
}
