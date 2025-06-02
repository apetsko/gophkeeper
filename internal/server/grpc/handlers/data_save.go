package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

		data, errMarshal := proto.Marshal(in.GetBankCard())
		if errMarshal != nil {
			return nil, fmt.Errorf("error serialize: %v", errMarshal)
		}

		encryptedData, errEncrypt := s.Envelop.EncryptUserData(
			ctx,
			encryptedMK,
			data,
		)
		if errEncrypt != nil {
			slog.Error("failed to crypt data: %w", errEncrypt)
			return nil, fmt.Errorf("failed to crypt data: %v", errEncrypt)
		}

		saveUserData := &models.DbUserData{
			UserID:        userID,
			Type:          constants.BankCard,
			EncryptedData: encryptedData.EncryptedData,
			DataNonce:     encryptedData.DataNonce,
			EncryptedDek:  encryptedData.EncryptedDek,
			DekNonce:      encryptedData.DekNonce,
			Meta:          protojson.Format(in.Meta),
		}
		_, errSave := s.Storage.SaveUserData(ctx, saveUserData)
		if errSave != nil {
			return nil, errSave
		}

	case pbc.DataType_DATA_TYPE_CREDENTIALS:
		creds := in.GetCredentials()
		if creds == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют учетные данные")
		}

		data, errMarshal := proto.Marshal(in.GetCredentials())
		if errMarshal != nil {
			return nil, fmt.Errorf("error serialize: %v", errMarshal)
		}

		encryptedData, errEncrypt := s.Envelop.EncryptUserData(
			ctx,
			encryptedMK,
			data,
		)
		if errEncrypt != nil {
			slog.Error("failed to crypt data: %w", errEncrypt)
			return nil, fmt.Errorf("failed to crypt data: %v", errEncrypt)
		}

		saveUserData := &models.DbUserData{
			UserID:        userID,
			Type:          constants.Credentials,
			EncryptedData: encryptedData.EncryptedData,
			DataNonce:     encryptedData.DataNonce,
			EncryptedDek:  encryptedData.EncryptedDek,
			DekNonce:      encryptedData.DekNonce,
			Meta:          protojson.Format(in.Meta),
		}
		_, errSave := s.Storage.SaveUserData(ctx, saveUserData)
		if errSave != nil {
			return nil, errSave
		}

	case pbc.DataType_DATA_TYPE_BINARY_DATA:
		file := in.GetBinaryData()
		if file == nil {
			return nil, status.Errorf(codes.InvalidArgument, "отсутствуют данные файла")
		}

		// Шифруем данные файла
		encryptedData, errEncrypt := s.Envelop.EncryptUserData(
			ctx,
			encryptedMK,
			file.Data,
		)
		if errEncrypt != nil {
			slog.Error("failed to encrypt binary data", "error", errEncrypt)
			return nil, fmt.Errorf("failed to encrypt binary data: %v", errEncrypt)
		}

		// Генерируем уникальное имя файла
		objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Name)

		s3UploadData := &models.S3UploadData{
			ObjectName:  objectName,
			MetaContent: in.Meta.Content,
			FileName:    file.Name,
			FileType:    file.Type,
		}

		_, errUpload := s.StorageS3.Upload(ctx, encryptedData.EncryptedData, s3UploadData)
		if errUpload != nil {
			return nil, fmt.Errorf("failed to upload file to MinIO: %v", errUpload)
		}

		saveUserData := &models.DbUserData{
			UserID:        userID,
			Type:          constants.BinaryData,
			MinioObjectID: objectName,
			DataNonce:     encryptedData.DataNonce,
			EncryptedDek:  encryptedData.EncryptedDek,
			DekNonce:      encryptedData.DekNonce,
			Meta:          protojson.Format(in.Meta),
		}
		_, errSave := s.Storage.SaveUserData(ctx, saveUserData)
		if errSave != nil {
			return nil, errSave
		}

	default:
		return nil, status.Errorf(codes.Unimplemented, "неподдерживаемый тип данных: %v", in.Type)
	}

	return &pbrpc.DataSaveResponse{
		Message: fmt.Sprintf("данные типа %s успешно сохранены", in.Type.String()),
	}, nil
}
