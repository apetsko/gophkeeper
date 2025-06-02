package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/models"
)

type S3 struct {
	MinioClient *minio.Client
	MinioBucket string
}

func NewS3Client(ctx context.Context, cfg config.MinioConfig) (*S3, error) {
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

	return &S3{
		MinioClient: minioClient,
		MinioBucket: cfg.Bucket,
	}, err
}

func (s *S3) Upload(
	ctx context.Context,
	data []byte,
	s3UploadData *models.S3UploadData,
) (*minio.UploadInfo, error) {
	info, errPutObject := s.MinioClient.PutObject(
		ctx,
		s.MinioBucket,
		s3UploadData.ObjectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: s3UploadData.FileType,
			UserMetadata: map[string]string{
				"original-name": s3UploadData.FileName,
				"meta-content":  s3UploadData.MetaContent,
				"upload-time":   time.Now().Format(time.RFC3339),
				"is-encrypted":  "true",
			},
		},
	)

	if errPutObject != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %v", errPutObject)
	}

	return &info, nil
}
