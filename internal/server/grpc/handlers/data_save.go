package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/models"
	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) DataSave(ctx context.Context, in *pbrpc.DataSaveRequest) (*pbrpc.DataSaveResponse, error) {
	userID, ok := ctx.Value(constants.UserID).(int)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось получить UserID")
	}

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

		// TODO: переделать на потокобезопасную in memory мапу
		encryptedMK, err := s.KeyManager.GetMasterKey(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("error get encryptedMK: %v", err)
		}

		data, err := proto.Marshal(in.GetBankCard())
		if err != nil {
			return nil, fmt.Errorf("error serialize: %v", err)
		}

		userData := &models.UserData{
			UserID:        userID,
			Type:          "bank_card", // TODO: типы унести в константы
			MinioObjectID: "",
			Meta:          protojson.Format(in.Meta),
		}

		_, errEncrypt := s.Envelop.EncryptUserData(
			ctx,
			*userData,
			encryptedMK,
			data,
		)
		if errEncrypt != nil {
			slog.Error("failed to save bank card: %w", errEncrypt)
			return nil, fmt.Errorf("failed to save bank card: %v", errEncrypt)
		}

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

		info, errPutObject := s.MinioClient.PutObject(
			ctx,
			s.MinioBucket,
			objectName,
			bytes.NewReader(file.Data),
			int64(len(file.Data)),
			minio.PutObjectOptions{
				ContentType: file.Type,
				UserMetadata: map[string]string{
					"original-name": file.Name,
					//"meta-content":  in.Meta.Content, panics
					"meta-content": "in.Meta.Content",
					"upload-time":  time.Now().Format(time.RFC3339),
				},
			},
		)

		//todo добавить сохранение инфы в базу, чтобы потом лист из базы делать, и потом при выборе файла прислать его.

		if errPutObject != nil {
			return nil, fmt.Errorf("failed to upload file to MinIO: %v", errPutObject)
		}

		log.Printf("Успешно загружен %s в бакет %s. ETAG: %s", objectName, s.MinioBucket, info.ETag)

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
