package stogage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/apetsko/gophkeeper/config"
)

func NewMinioClient(ctx context.Context, cfg config.MinioConfig) (*minio.Client, error) {
	var err error

	minioClient, err := minio.New(cfg.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.ID, cfg.Secret, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("error init minio client: %v", err)
	}

	bucketExists, err := minioClient.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !bucketExists {
		if errMakeBucket := minioClient.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); errMakeBucket != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", errMakeBucket)
		}
	}

	return minioClient, err
}
