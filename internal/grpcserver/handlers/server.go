package handlers

import "github.com/minio/minio-go/v7"

type ServerAdmin struct {
	minioBucket string
	minioClient *minio.Client
}

func NewServer(
	minioBucket string,
	minioClient *minio.Client,
) *ServerAdmin {
	return &ServerAdmin{
		minioBucket: minioBucket,
		minioClient: minioClient,
	}
}
