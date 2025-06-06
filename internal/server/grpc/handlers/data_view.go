// Package handlers provides gRPC server handlers for managing user data operations,
// including creation, retrieval, update, and deletion of user records.
package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/apetsko/gophkeeper/internal/constants"
	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	"github.com/apetsko/gophkeeper/protogen/api/proto/v1/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

var stringToDataType = map[string]pbc.DataType{
	"bank_card":   pbc.DataType_DATA_TYPE_BANK_CARD,
	"credentials": pbc.DataType_DATA_TYPE_CREDENTIALS,
	"binary_data": pbc.DataType_DATA_TYPE_BINARY_DATA,
}

// DataView handles the gRPC request to retrieve a specific user data record by its ID.
//
// This method checks user authorization, fetches the encrypted data from the database or S3 (for binary files),
// decrypts the data using the user's master key, parses it according to its type (bank card, credentials, or binary data),
// and returns the result in the response.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The DataViewRequest message containing the record ID.
//
// Returns:
//   - *pbrpc.DataViewResponse: The requested user data record.
//   - error: A gRPC error if access is denied or an internal error occurs.
func (s *ServerAdmin) DataView(ctx context.Context, in *pbrpc.DataViewRequest) (*pbrpc.DataViewResponse, error) {
	userID, ok := ctx.Value(constants.UserID).(int)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось получить UserID")
	}

	userData, err := s.Storage.GetUserData(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных")
	}

	// Проверка прав доступа
	if userData.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "нет доступа к запрошенным данным")
	}

	encryptedMK, err := s.KeyManager.GetMasterKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error get encryptedMK: %v", err)
	}

	var decryptData []byte
	var dataType pbc.DataType
	var file models.File

	switch userData.Type {
	case constants.BinaryData:
		// Обработка бинарных данных (скачивание из MinIO)
		fileData, fileInfo, errGetObject := s.StorageS3.GetObject(
			ctx,
			userData.MinioObjectID,
		)
		if errGetObject != nil {
			return nil, status.Errorf(codes.Internal, "ошибка получения файла из хранилища: %v", errGetObject)
		}

		userData.EncryptedData = fileData

		decryptData, err = s.Envelope.DecryptUserData(ctx, *userData, encryptedMK)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ошибка расшифровки файла: %v", err)
		}

		originalName := fileInfo.UserMetadata["original-name"]
		if originalName == "" {
			originalName = userData.MinioObjectID
		}

		file = models.File{
			Name: originalName,
			Data: decryptData,
			Size: int32(len(decryptData)),
			Type: fileInfo.ContentType,
		}

		dataType = pbc.DataType_DATA_TYPE_BINARY_DATA
	default:
		decryptData, err = s.Envelope.DecryptUserData(ctx, *userData, encryptedMK)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ошибка расшифровки данных")
		}

		var ok bool
		dataType, ok = stringToDataType[userData.Type]
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "неподдерживаемый тип данных: %s", userData.Type)
		}
	}

	var meta models.Meta
	if errUnmarshal := protojson.Unmarshal([]byte(userData.Meta), &meta); errUnmarshal != nil {
		return nil, status.Errorf(codes.Internal, "ошибка парсинга Meta JSON: %v", errUnmarshal)
	}

	// 4. Создаем базовый ответ
	response := &pbrpc.DataViewResponse{
		Type: dataType,
		Meta: &meta,
	}

	// 5. Парсим данные в зависимости от типа
	if err := parseData(response, dataType, decryptData, &file); err != nil {
		return nil, err
	}
	return response, nil
}

func parseData(r *pbrpc.DataViewResponse, dataType pbc.DataType, decryptData []byte, file *models.File) error {
	switch dataType {
	case pbc.DataType_DATA_TYPE_BANK_CARD:
		var card models.BankCard
		if errUnmarshal := proto.Unmarshal(decryptData, &card); errUnmarshal != nil {
			return status.Errorf(codes.Internal, "ошибка парсинга карты")
		}
		r.Data = &pbrpc.DataViewResponse_BankCard{BankCard: &card}

	case pbc.DataType_DATA_TYPE_CREDENTIALS:
		var credentials models.Credentials
		if errUnmarshal := proto.Unmarshal(decryptData, &credentials); errUnmarshal != nil {
			return status.Errorf(codes.Internal, "ошибка парсинга учетных данных")
		}
		r.Data = &pbrpc.DataViewResponse_Credentials{Credentials: &credentials}

	case pbc.DataType_DATA_TYPE_BINARY_DATA:
		r.Data = &pbrpc.DataViewResponse_BinaryData{BinaryData: file}

	default:
		return status.Errorf(codes.InvalidArgument, "неподдерживаемый тип данных")
	}
	return nil
}
