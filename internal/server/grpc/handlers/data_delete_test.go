package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerAdmin_DataDelete(t *testing.T) {
	const userID = 42
	const dataID = 101

	tests := []struct {
		name           string
		ctx            context.Context
		req            *pbrpc.DataDeleteRequest
		mockSetup      func(m *mocks.IStorage)
		wantErr        bool
		wantErrCode    codes.Code
		wantErrMessage string
		wantRespMsg    string
	}{
		{
			name:           "missing userID in context",
			ctx:            context.Background(),
			req:            &pbrpc.DataDeleteRequest{Id: int32(dataID)},
			mockSetup:      func(m *mocks.IStorage) {},
			wantErr:        true,
			wantErrCode:    codes.InvalidArgument,
			wantErrMessage: "не удалось получить UserID",
		},
		{
			name: "GetUserData returns error",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			req:  &pbrpc.DataDeleteRequest{Id: int32(dataID)},
			mockSetup: func(m *mocks.IStorage) {
				m.On("GetUserData", mock.Anything, dataID).
					Return(nil, errors.New("db error")).Once()
			},
			wantErr:        true,
			wantErrCode:    codes.Internal,
			wantErrMessage: "ошибка получения данных",
		},
		{
			name: "userID mismatch",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			req:  &pbrpc.DataDeleteRequest{Id: int32(dataID)},
			mockSetup: func(m *mocks.IStorage) {
				m.On("GetUserData", mock.Anything, dataID).
					Return(&models.DBUserData{UserID: 999}, nil).Once()
			},
			wantErr:        true,
			wantErrCode:    codes.PermissionDenied,
			wantErrMessage: "нельзя удалить запись, она не ваша",
		},
		{
			name: "DeleteUserData returns error",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			req:  &pbrpc.DataDeleteRequest{Id: int32(dataID)},
			mockSetup: func(m *mocks.IStorage) {
				m.On("GetUserData", mock.Anything, dataID).
					Return(&models.DBUserData{UserID: userID}, nil).Once()
				m.On("DeleteUserData", mock.Anything, dataID).
					Return(errors.New("delete error")).Once()
			},
			wantErr:        true,
			wantErrCode:    codes.Internal,
			wantErrMessage: "ошибка удаления данных",
		},
		{
			name: "successful delete",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			req:  &pbrpc.DataDeleteRequest{Id: int32(dataID)},
			mockSetup: func(m *mocks.IStorage) {
				m.On("GetUserData", mock.Anything, dataID).
					Return(&models.DBUserData{UserID: userID}, nil).Once()
				m.On("DeleteUserData", mock.Anything, dataID).
					Return(nil).Once()
			},
			wantErr:     false,
			wantRespMsg: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &mocks.IStorage{}
			tt.mockSetup(mockStorage)

			srv := &ServerAdmin{
				Storage: mockStorage,
			}

			resp, err := srv.DataDelete(tt.ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantErrMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantRespMsg, resp.GetMessage())
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
