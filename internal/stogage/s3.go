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

		if errSetBucketPolicy := SetBucketPolicy(ctx, minioClient, cfg.Bucket); errSetBucketPolicy != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %v", errSetBucketPolicy)
		}
	}

	return minioClient, err
}

func SetBucketPolicy(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::gophkeeper/*"],
				"Condition": {}
			},
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:ListBucket"],
				"Resource": ["arn:aws:s3:::gophkeeper"],
				"Condition": {}
			}
		]
	}`

	return minioClient.SetBucketPolicy(ctx, bucketName, policy)
}
