package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/models"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
)

func getTestS3Config() config.S3Config {
	return config.S3Config{
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "gophkeeper",
		Endpoint:  "localhost:9000",
	}
}

func TestS3_FullFlow(t *testing.T) {
	cfg := getTestS3Config()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s3, err := NewS3Client(ctx, cfg)
	require.NoError(t, err)

	objectName := "test-object.txt"
	content := []byte("hello minio test")
	uploadData := &models.S3UploadData{
		ObjectName:  objectName,
		FileName:    "test-object.txt",
		FileType:    "text/plain",
		MetaContent: "test-meta",
	}

	// Upload
	uploadInfo, err := s3.Upload(ctx, content, uploadData)
	require.NoError(t, err)
	require.Equal(t, objectName, uploadInfo.Key)

	// Download
	got, info, err := s3.GetObject(ctx, objectName)
	require.NoError(t, err)
	fmt.Printf("%#v", info)
	require.Equal(t, int64(len(content)), info.Size)
	require.Equal(t, content, got)
	require.Equal(t, "text/plain", info.ContentType)
	require.Equal(t, "test-object.txt", info.Key)
	require.Equal(t, "test-meta", info.UserMetadata["Meta-Content"])
	require.Equal(t, "true", info.UserMetadata["Is-Encrypted"])

	// Delete
	err = s3.MinioClient.RemoveObject(ctx, cfg.Bucket, objectName, minio.RemoveObjectOptions{})
	require.NoError(t, err)

	// Ensure deleted
	_, _, err = s3.GetObject(ctx, objectName)
	require.Error(t, err)
}
