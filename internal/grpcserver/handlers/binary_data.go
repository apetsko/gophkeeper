package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"

	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) BinaryData(ctx context.Context, in *pbrpc.BinaryDataRequest) (*pbrpc.BinaryDataResponse, error) {
	if in.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), in.File.Name) // Добавляем timestamp для уникальности

	// Загружаем файл в MinIO
	info, errPutObject := s.minioClient.PutObject(
		ctx,
		s.minioBucket,
		objectName,
		bytes.NewReader(in.File.Data),
		int64(len(in.File.Data)),
		minio.PutObjectOptions{
			ContentType: in.File.Type,
			UserMetadata: map[string]string{
				"original-name": in.File.Name,
				"meta-content":  in.Meta.Content,
				"upload-time":   time.Now().Format(time.RFC3339),
			},
		},
	)

	if errPutObject != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %v", errPutObject)
	}

	log.Printf("Successfully uploaded %s to bucket %s. ETAG: %s", objectName, s.minioBucket, info.ETag)

	// todo: здесь нужно сохранить метаданные и ссылку на файл в БД
	if in.Meta != nil {
		log.Printf("Received metadata for file: %+v", in.Meta)
	}

	return &pbrpc.BinaryDataResponse{
		Message: fmt.Sprintf("File uploaded successfully. Object name: %s", objectName),
	}, nil
}
