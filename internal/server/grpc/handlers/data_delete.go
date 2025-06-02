package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/apetsko/gophkeeper/internal/constants"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) DataDelete(ctx context.Context, in *pbrpc.DataDeleteRequest) (*pbrpc.DataDeleteResponse, error) {
	userID, ok := ctx.Value(constants.UserID).(int)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось получить UserID")
	}

	userData, err := s.Storage.GetUserData(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных")
	}

	if userData.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "нельзя удалить запись, она не ваша")
	}

	errDelete := s.Storage.DeleteUserData(ctx, int(in.GetId()))
	if errDelete != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных")
	}

	return &pbrpc.DataDeleteResponse{
		Message: "ok",
	}, nil
}
