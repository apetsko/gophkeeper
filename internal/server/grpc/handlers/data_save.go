// Package handlers provides gRPC server handlers for managing user data operations,
// including creation, retrieval, update, and deletion of user records.
package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/models"
	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	pbmodels "github.com/apetsko/gophkeeper/protogen/api/proto/v1/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DataSave handles the gRPC request to save user data.
//
// This method validates the request, retrieves the user's master key, encrypts the data,
// and stores it in the database or S3 depending on the data type.
//
// Parameters:
// - ctx: The gRPC context.
// - in: The DataSaveRequest message containing user data.
//
// Returns:
// - *pbrpc.DataSaveResponse: A response indicating success.
// - error: An error if validation or storage fails.
func (s *ServerAdmin) DataSave(ctx context.Context, in *pbrpc.DataSaveRequest) (*pbrpc.DataSaveResponse, error) {
	userID, ok := ctx.Value(constants.UserID).(int)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось получить UserID")
	}

	// Валидация типа данных
	if in.Type == pbc.DataType_DATA_TYPE_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "тип данных не указан")
	}

	// TODO: переделать на потокобезопасную in memory мапу
	encryptedMK, err := s.KeyManager.GetMasterKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error get encryptedMK: %v", err)
	}

	// Обработка данных в зависимости от типа
	switch in.Type {
	case pbc.DataType_DATA_TYPE_BANK_CARD:
		bankCard := in.GetBankCard()
		if bankCard == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют данные банковской карты")
		}

		err := s.saveUserData(ctx, userID, in.Type, encryptedMK, bankCard, in.Meta)
		if err != nil {
			return nil, err
		}

	case pbc.DataType_DATA_TYPE_CREDENTIALS:
		creds := in.GetCredentials()
		if creds == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют учетные данные")
		}
		err := s.saveUserData(ctx, userID, in.Type, encryptedMK, creds, in.Meta)
		if err != nil {
			return nil, err
		}

	case pbc.DataType_DATA_TYPE_BINARY_DATA:
		file := in.GetBinaryData()
		if file == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют данные файла")
		}

		// Шифруем содержимое файла
		encryptedData, err := s.Envelope.EncryptUserData(ctx, encryptedMK, file.Data)
		if err != nil {
			slog.Error("failed to encrypt binary data", "error", err)
			return nil, fmt.Errorf("failed to encrypt binary data: %v", err)
		}

		// Генерируем уникальное имя файла
		objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Name)

		// Загружаем в S3
		s3UploadData := &models.S3UploadData{
			ObjectName:  objectName,
			MetaContent: in.Meta.Content,
			FileName:    file.Name,
			FileType:    file.Type,
		}
		_, err = s.StorageS3.Upload(ctx, encryptedData.EncryptedData, s3UploadData)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file to MinIO: %v", err)
		}

		// Сохраняем метаданные в БД
		saveUserData := &models.DBUserData{
			UserID:        userID,
			Type:          constants.MapDataTypeToString(in.Type),
			MinioObjectID: objectName,
			DataNonce:     encryptedData.DataNonce,
			EncryptedDek:  encryptedData.EncryptedDek,
			DekNonce:      encryptedData.DekNonce,
			Meta:          protojson.Format(in.Meta),
		}
		_, err = s.Storage.SaveUserData(ctx, saveUserData)
		if err != nil {
			return nil, err
		}

	default:
		return nil, status.Errorf(codes.Unimplemented, "неподдерживаемый тип данных: %v", in.Type)
	}

	return &pbrpc.DataSaveResponse{
		Message: fmt.Sprintf("данные типа %s успешно сохранены", in.Type.String()),
	}, nil
}

func (s *ServerAdmin) saveUserData(
	ctx context.Context,
	userID int,
	dataType pbc.DataType,
	encryptedMK []byte,
	data proto.Message,
	meta *pbmodels.Meta,
) error {
	// Маршал protobuf
	serialized, err := proto.Marshal(data)
	if err != nil {
		return fmt.Errorf("serialize error: %v", err)
	}

	// Шифруем данные
	encryptedData, err := s.Envelope.EncryptUserData(ctx, encryptedMK, serialized)
	if err != nil {
		slog.Error("failed to crypt data: " + err.Error())
		return fmt.Errorf("encrypt error: %v", err)
	}

	// Сохраняем в БД
	saveUserData := &models.DBUserData{
		UserID:        userID,
		Type:          constants.MapDataTypeToString(dataType),
		EncryptedData: encryptedData.EncryptedData,
		DataNonce:     encryptedData.DataNonce,
		EncryptedDek:  encryptedData.EncryptedDek,
		DekNonce:      encryptedData.DekNonce,
		Meta:          protojson.Format(meta),
	}
	_, err = s.Storage.SaveUserData(ctx, saveUserData)
	return err
}
