// Package storage provides S3-compatible object storage integration for GophKeeper.
//
// This package implements an S3 client using MinIO for uploading and retrieving encrypted user files.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/models"
)

// S3Client defines the interface for S3-compatible object storage operations.
//
// It abstracts file upload and retrieval for easier testing and mocking.
type S3Client interface {
	Upload(ctx context.Context, data []byte, s3UploadData *models.S3UploadData) (*minio.UploadInfo, error)
	GetObject(ctx context.Context, objectName string) ([]byte, *minio.ObjectInfo, error)
}

// S3 implements the S3Client interface using a MinIO client.
//
// It provides methods to upload and retrieve objects from the configured bucket.
type S3 struct {
	MinioClient *minio.Client
	MinioBucket string
}

// NewS3Client initializes a new S3 client with the given configuration.
//
// It connects to the MinIO server, checks for the existence of the bucket, and creates it if necessary.
//
// Parameters:
//   - ctx: Context for the operation.
//   - cfg: S3Config with endpoint, credentials, and bucket name.
//
// Returns:
//   - *S3: The initialized S3 client.
//   - error: An error if initialization fails.
func NewS3Client(ctx context.Context, cfg config.S3Config) (*S3, error) {
	var err error

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
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

// Upload uploads data to the S3 bucket with the specified metadata.
//
// Parameters:
//   - ctx: Context for the operation.
//   - data: File data to upload.
//   - s3UploadData: Metadata and object information.
//
// Returns:
//   - *minio.UploadInfo: Information about the uploaded object.
//   - error: An error if the upload fails.
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

// GetObject retrieves an object from the S3 bucket by its name.
//
// Parameters:
//   - ctx: Context for the operation.
//   - objectName: Name of the object to retrieve.
//
// Returns:
//   - []byte: The object data.
//   - *minio.ObjectInfo: Metadata about the object.
//   - error: An error if retrieval fails.
func (s *S3) GetObject(
	ctx context.Context,
	objectName string,
) ([]byte, *minio.ObjectInfo, error) {
	object, err := s.MinioClient.GetObject(
		ctx,
		s.MinioBucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get object from MinIO: %v", err)
	}
	defer object.Close()

	objectInfo, err := object.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get object info: %v", err)
	}

	data := make([]byte, objectInfo.Size)
	_, err = io.ReadFull(object, data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read object data: %v", err)
	}

	return data, &objectInfo, nil
}
