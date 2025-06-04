package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/models"

	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerAdmin_DataList(t *testing.T) {
	const userID = 42

	tests := []struct {
		name            string
		ctx             context.Context
		mockSetup       func(m *mocks.IStorage)
		wantErr         bool
		wantErrCode     codes.Code
		wantErrContains string
		wantCount       int32
		wantRecordsLen  int
	}{
		{
			name:            "missing userID in context",
			ctx:             context.Background(),
			mockSetup:       func(m *mocks.IStorage) {},
			wantErr:         true,
			wantErrCode:     codes.InvalidArgument,
			wantErrContains: "не удалось получить UserID",
		},
		{
			name: "storage returns error",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			mockSetup: func(m *mocks.IStorage) {
				m.On("GetUserDataList", mock.Anything, userID).
					Return(nil, errors.New("db error")).Once()
			},
			wantErr:         true,
			wantErrCode:     codes.Internal,
			wantErrContains: "ошибка получения данных",
		},
		{
			name: "successful list with valid meta",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			mockSetup: func(m *mocks.IStorage) {
				now := time.Now()
				m.On("GetUserDataList", mock.Anything, userID).
					Return([]models.UserDataListItem{
						{
							ID:        1,
							UserID:    userID,
							Type:      "note",
							Meta:      `{"content":"test content"}`,
							CreatedAt: now,
						},
						{
							ID:        2,
							UserID:    userID,
							Type:      "card",
							Meta:      `{"content":"another content"}`,
							CreatedAt: now.Add(time.Minute),
						},
					}, nil).Once()
			},
			wantErr:        false,
			wantCount:      2,
			wantRecordsLen: 2,
		},
		{
			name: "one record with invalid meta JSON skipped",
			ctx:  context.WithValue(context.Background(), constants.UserID, userID),
			mockSetup: func(m *mocks.IStorage) {
				now := time.Now()
				m.On("GetUserDataList", mock.Anything, userID).
					Return([]models.UserDataListItem{
						{
							ID:        1,
							UserID:    userID,
							Type:      "note",
							Meta:      `{"content":"valid"}`,
							CreatedAt: now,
						},
						{
							ID:        2,
							UserID:    userID,
							Type:      "card",
							Meta:      `invalid json`,
							CreatedAt: now.Add(time.Minute),
						},
					}, nil).Once()
			},
			wantErr:        false,
			wantCount:      1,
			wantRecordsLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &mocks.IStorage{}
			tt.mockSetup(mockStorage)

			srv := &ServerAdmin{
				Storage: mockStorage,
			}

			resp, err := srv.DataList(tt.ctx, &pbrpc.DataListRequest{})

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantErrContains)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantCount, resp.GetCount())
				assert.Len(t, resp.GetRecords(), tt.wantRecordsLen)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
