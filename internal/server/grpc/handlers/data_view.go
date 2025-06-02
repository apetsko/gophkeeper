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

func (s *ServerAdmin) DataView(ctx context.Context, in *pbrpc.DataViewRequest) (*pbrpc.DataViewResponse, error) {
	userID, ok := ctx.Value(constants.UserID).(int)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось получить UserID")
	}

	userData, err := s.Storage.GetUserData(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных")
	}

	// TODO: переделать на потокобезопасную in memory мапу
	encryptedMK, err := s.KeyManager.GetMasterKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error get encryptedMK: %v", err)
	}

	// TODO: нет обработчика для файла

	decryptData, err := s.Envelop.DecryptUserData(ctx, *userData, encryptedMK)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка расшифровки данных")
	}

	dataType, ok := stringToDataType[userData.Type]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "неподдерживаемый тип данных: %s", userData.Type)
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
	switch dataType {
	case pbc.DataType_DATA_TYPE_BANK_CARD:
		var card models.BankCard
		if errUnmarshal := proto.Unmarshal(decryptData, &card); errUnmarshal != nil {
			return nil, status.Errorf(codes.Internal, "ошибка парсинга карты")
		}
		response.Data = &pbrpc.DataViewResponse_BankCard{BankCard: &card}

	case pbc.DataType_DATA_TYPE_CREDENTIALS:
		var creds models.Credentials
		if errUnmarshal := proto.Unmarshal(decryptData, &creds); errUnmarshal != nil {
			return nil, status.Errorf(codes.Internal, "ошибка парсинга учетных данных")
		}
		response.Data = &pbrpc.DataViewResponse_Credentials{Credentials: &creds}

	case pbc.DataType_DATA_TYPE_BINARY_DATA:
		var file models.File
		if errUnmarshal := proto.Unmarshal(decryptData, &file); errUnmarshal != nil {
			return nil, status.Errorf(codes.Internal, "ошибка парсинга файла")
		}
		response.Data = &pbrpc.DataViewResponse_BinaryData{BinaryData: &file}

	default:
		return nil, status.Errorf(codes.InvalidArgument, "неподдерживаемый тип данных")
	}

	return response, nil
}
