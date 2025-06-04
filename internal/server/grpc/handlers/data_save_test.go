package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/gophkeeper/internal/constants"

	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/models"
	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	pbmodels "github.com/apetsko/gophkeeper/protogen/api/proto/v1/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerAdmin_DataSave(t *testing.T) {
	userID := 42
	ctx := context.WithValue(context.Background(), constants.UserID, userID)
	encryptedMK := []byte("masterkey")

	encryptMockData := &models.EncryptedData{
		EncryptedData: []byte("encryptedData"),
		DataNonce:     []byte("nonce1"),
		EncryptedDek:  []byte("encryptedDEK"),
		DekNonce:      []byte("nonce2"),
	}

	tests := []struct {
		name          string
		ctx           context.Context
		req           *pbrpc.DataSaveRequest
		mockSetup     func(s *ServerAdmin)
		wantErr       bool
		wantErrCode   codes.Code
		wantMsgSubstr string
	}{
		{
			name: "no userID in context",
			ctx:  context.Background(),
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Data: &pbrpc.DataSaveRequest_BankCard{
					BankCard: &pbmodels.BankCard{
						CardNumber: "1111222233334444",
						ExpiryDate: "12/25",
						Cvv:        "123",
						Cardholder: "John Doe",
					},
				},
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "type unspecified",
			ctx:  ctx,
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_UNSPECIFIED,
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "keymanager returns error",
			ctx:  ctx,
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_CREDENTIALS,
				Data: &pbrpc.DataSaveRequest_Credentials{
					Credentials: &pbmodels.Credentials{
						Login:    "user",
						Password: "pass",
					},
				},
			},
			mockSetup: func(s *ServerAdmin) {
				s.KeyManager = new(mocks.KeyManagerInterface)
				s.KeyManager.On("GetMasterKey", ctx, userID).Return(nil, errors.New("keymgr error"))
			},
			wantErr: true,
		},
		{
			name: "bank card with nil data",
			ctx:  ctx,
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Data: nil,
			},
			mockSetup: func(s *ServerAdmin) {
				s.KeyManager = new(mocks.KeyManagerInterface)
				s.KeyManager.On("GetMasterKey", ctx, userID).Return(encryptedMK, nil)
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "successful save bank card",
			ctx:  ctx,
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Meta: &pbmodels.Meta{Content: "meta info"},
				Data: &pbrpc.DataSaveRequest_BankCard{
					BankCard: &pbmodels.BankCard{
						CardNumber: "1234567812345678",
						ExpiryDate: "01/30",
						Cvv:        "321",
						Cardholder: "Alice",
					},
				},
			},
			mockSetup: func(s *ServerAdmin) {
				s.KeyManager = new(mocks.KeyManagerInterface)
				s.KeyManager.On("GetMasterKey", ctx, userID).Return(encryptedMK, nil)

				s.Envelope = new(mocks.IEnvelope)
				s.Envelope.On("EncryptUserData", ctx, encryptedMK, mock.Anything).Return(encryptMockData, nil)

				s.Storage = new(mocks.IStorage)
				s.Storage.On("SaveUserData", ctx, mock.AnythingOfType("*models.DBUserData")).Return(1, nil)
			},
			wantErr: false,
		},
		{
			name: "successful save binary data",
			ctx:  ctx,
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BINARY_DATA,
				Meta: &pbmodels.Meta{Content: "file meta"},
				Data: &pbrpc.DataSaveRequest_BinaryData{
					BinaryData: &pbmodels.File{
						Name: "file.txt",
						Type: "text/plain",
						Size: 10,
						Data: []byte("hello world"),
					},
				},
			},
			mockSetup: func(s *ServerAdmin) {
				s.KeyManager = new(mocks.KeyManagerInterface)
				s.KeyManager.On("GetMasterKey", ctx, userID).Return(encryptedMK, nil)

				s.Envelope = new(mocks.IEnvelope)
				s.Envelope.On("EncryptUserData", ctx, encryptedMK, []byte("hello world")).Return(encryptMockData, nil)

				s.StorageS3 = new(mocks.S3Client)
				s.StorageS3.On("Upload", ctx, encryptMockData.EncryptedData, mock.AnythingOfType("*models.S3UploadData")).Return("objName", nil)

				s.Storage = new(mocks.IStorage)
				s.Storage.On("SaveUserData", ctx, mock.AnythingOfType("*models.DBUserData")).Return(1, nil)
			},
			wantErr: false,
		},
		{
			name: "unsupported data type",
			ctx:  ctx,
			req: &pbrpc.DataSaveRequest{
				Type: 999,
				Data: nil,
			},
			mockSetup: func(s *handlers.ServerAdmin) {
				s.KeyManager = new(mock.KeyManager)
			},
			wantErr:     true,
			wantErrCode: codes.Unimplemented,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &handlers.ServerAdmin{}
			if tt.mockSetup != nil {
				tt.mockSetup(srv)
			}

			resp, err := srv.DataSave(tt.ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Message, tt.req.Type.String())
			}
		})
	}
}
