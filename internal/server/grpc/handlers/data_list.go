// Package handlers provides gRPC server handlers for managing user data operations,
// including creation, retrieval, update, and deletion of user records.
package handlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/apetsko/gophkeeper/internal/constants"
	pbmodels "github.com/apetsko/gophkeeper/protogen/api/proto/v1/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

// DataView handles the gRPC request to retrieve a specific user data record by its ID.
//
// This method checks user authorization, fetches and decrypts the data (including files from S3 if needed),
// parses the data according to its type, and returns it in the response.
//
// Parameters:
// - ctx: The gRPC context.
// - in: The DataViewRequest message with the data record ID.
//
// Returns:
// - *pbrpc.DataViewResponse: The requested data record.
// - error: An error if access is denied or retrieval fails.
func (s *ServerAdmin) DataList(ctx context.Context, in *pbrpc.DataListRequest) (*pbrpc.DataListResponse, error) {
	userID, ok := ctx.Value(constants.UserID).(int)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось получить UserID")
	}

	userDataList, err := s.Storage.GetUserDataList(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных: %v", err)
	}

	var records []*pbmodels.Record
	for _, data := range userDataList {
		// Преобразуем строку Meta (в формате JSON) в pbmodels.Meta
		var meta pbmodels.Meta
		if errUnmarshal := protojson.Unmarshal([]byte(data.Meta), &meta); errUnmarshal != nil {
			slog.Error("failed to unmarshal meta: " + errUnmarshal.Error())
			continue
		}

		record := &pbmodels.Record{
			Id:        int32(data.ID),
			Type:      data.Type,
			Meta:      &meta,
			CreatedAt: data.CreatedAt.Format("02.01.2006 15:04"),
		}
		records = append(records, record)
	}

	return &pbrpc.DataListResponse{
		Records: records,
		Count:   int32(len(records)), // Или общее количество в базе, если есть пагинация
	}, nil
}
