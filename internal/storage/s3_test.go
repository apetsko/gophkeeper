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
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var minioContainer tc.Container

func startTestMinio(t *testing.T) (endpoint string, terminate func(), err error) {
	ctx := context.Background()
	req := tc.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForListeningPort("9000/tcp").WithStartupTimeout(20 * time.Second),
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if t != nil {
			require.NoError(t, err)
		}
		return "", nil, err
	}
	minioContainer = container

	host, err := container.Host(ctx)
	if err != nil {
		if t != nil {
			require.NoError(t, err)
		}
		return "", nil, err
	}
	port, err := container.MappedPort(ctx, "9000")
	if err != nil {
		if t != nil {
			require.NoError(t, err)
		}
		return "", nil, err
	}
	endpoint = fmt.Sprintf("%s:%s", host, port.Port())

	return endpoint, func() { _ = container.Terminate(ctx) }, nil
}

func TestS3_FullFlow(t *testing.T) {
	endpoint, terminate, err := startTestMinio(t)
	require.NoError(t, err)
	defer terminate()

	cfg := config.S3Config{
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "gophkeeper",
		Endpoint:  endpoint,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	// Retry upload if MinIO is not ready
	var uploadInfo *minio.UploadInfo
	maxUploadAttempts := 10
	for i := 0; i < maxUploadAttempts; i++ {
		uploadInfo, err = s3.Upload(ctx, content, uploadData)
		if err == nil {
			break
		}
		if i == maxUploadAttempts-1 {
			require.NoError(t, err, "failed to upload after retries")
		}
		time.Sleep(1 * time.Second)
	}
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
