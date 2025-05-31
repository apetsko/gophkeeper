package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"

	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAdmin) DataSave(ctx context.Context, in *pbrpc.DataSaveRequest) (*pbrpc.DataSaveResponse, error) {
	// Валидация типа данных
	if in.Type == pbc.DataType_DATA_TYPE_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "тип данных не указан")
	}

	// Обработка данных в зависимости от типа
	switch in.Type {
	case pbc.DataType_DATA_TYPE_BANK_CARD:
		bankCard := in.GetBankCard()
		if bankCard == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют данные банковской карты")
		}
		fmt.Printf("Сохранение карты: %s (%s)\n", bankCard.GetCardNumber(), bankCard.GetCardholder())

	case pbc.DataType_DATA_TYPE_CREDENTIALS:
		creds := in.GetCredentials()
		if creds == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют учетные данные")
		}
		fmt.Printf("Сохранение учётных данных: %s\n", creds.GetLogin())

	case pbc.DataType_DATA_TYPE_BINARY_DATA:
		file := in.GetBinaryData()
		if file == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют данные файла")
		}

		objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Name)

		fmt.Printf("Сохранение файла: %s (%d bytes)\n", file.Name, file.Size)

		info, errPutObject := s.minioClient.PutObject(
			ctx,
			s.minioBucket,
			objectName,
			bytes.NewReader(file.Data),
			int64(len(file.Data)),
			minio.PutObjectOptions{
				ContentType: file.Type,
				UserMetadata: map[string]string{
					"original-name": file.Name,
					"meta-content":  in.Meta.Content,
					"upload-time":   time.Now().Format(time.RFC3339),
				},
			},
		)

		if errPutObject != nil {
			return nil, fmt.Errorf("failed to upload file to MinIO: %v", errPutObject)
		}

		log.Printf("Успешно загружен %s в бакет %s. ETAG: %s", objectName, s.minioBucket, info.ETag)

	default:
		return nil, status.Errorf(codes.Unimplemented, "неподдерживаемый тип данных: %v", in.Type)
	}

	if in.Meta != nil {
		fmt.Printf("Метаданные: %+v\n", in.Meta)
	}

	return &pbrpc.DataSaveResponse{
		Message: fmt.Sprintf("данные типа %s успешно сохранены", in.Type.String()),
	}, nil
}
