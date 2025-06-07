package storage

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/models"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
)

const (
	minioContainerName = "test_minio_container"
	minioEndpoint      = "localhost:9000"
)

func startTestMinio() {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"--name", minioContainerName,
		"-e", "MINIO_ROOT_USER=minioadmin",
		"-e", "MINIO_ROOT_PASSWORD=minioadmin",
		"-p", "9000:9000",
		"minio/minio", "server", "/data")
	_ = cmd.Run()
	waitForMinio()
}

func stopTestMinio() {
	_ = exec.Command("docker", "stop", minioContainerName).Run()
}

func waitForMinio() {
	timeout := time.After(20 * time.Second)
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-timeout:
			os.Exit(1)
		case <-tick:
			conn, err := net.DialTimeout("tcp", minioEndpoint, 1*time.Second)
			if err == nil {
				conn.Close()
				return
			}
		}
	}
}

func waitForMinioReady() {
	url := "http://localhost:9000/minio/health/ready"
	if isCI {
		url = "http://minio:9000/minio/health/ready"
	}
	timeout := time.After(60 * time.Second)
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			panic("Timeout waiting for MinIO to be ready")
		case <-tick:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

func getTestS3Config() config.S3Config {
	return config.S3Config{
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "gophkeeper",
		Endpoint:  "localhost:9000",
	}
}

func TestS3_FullFlow(t *testing.T) {
	//TODO fix workflow

	cfg := getTestS3Config()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	s3, err := NewS3Client(ctx, cfg)
	require.NoError(t, err)

	// Wait for bucket to be ready before upload
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		exists, err := s3.MinioClient.BucketExists(ctx, cfg.Bucket)
		if err == nil && exists {
			fmt.Println("Bucket exists")
			time.Sleep(2 * time.Second)
			break
		}
		time.Sleep(2 * time.Second)
		if i == maxAttempts-1 {
			require.NoError(t, fmt.Errorf("bucket not ready after retries: %v", err))
		}
	}

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
