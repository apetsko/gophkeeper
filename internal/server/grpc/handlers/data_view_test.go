package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/internal/mocks"

	"github.com/apetsko/gophkeeper/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServerAdmin_DataView(t *testing.T) {
	const userID = 42
	ctx := context.WithValue(context.Background(), constants.UserID, userID)

	type testCase struct {
		name       string
		req        *pbrpc.DataViewRequest
		setupMocks func(
			st *mocks.IStorage,
			s3 *mocks.S3Client,
			env *mocks.IEnvelope,
			km *mocks.KeyManagerInterface,
		)
		wantErr bool
	}

	tests := []testCase{
		{
			name: "success bank card",
			req:  &pbrpc.DataViewRequest{Id: 1},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 1).Return(&models.DBUserData{
					UserID: userID,
					Type:   "bank_card",
					Meta:   `{"content":"meta"}`,
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("DecryptUserData", mock.Anything, mock.AnythingOfType("models.DBUserData"), []byte("mk")).
					Return([]byte{10, 4, '1', '2', '3', '4'}, nil) // serialized BankCard
			},
			wantErr: false,
		},
		{
			name: "success credentials",
			req:  &pbrpc.DataViewRequest{Id: 2},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 2).Return(&models.DBUserData{
					UserID: userID,
					Type:   "credentials",
					Meta:   `{"content":"meta"}`,
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("DecryptUserData", mock.Anything, mock.AnythingOfType("models.DBUserData"), []byte("mk")).
					Return([]byte{10, 1, 'l', 18, 1, 'p'}, nil) // serialized Credentials
			},
			wantErr: false,
		},
		{
			name: "success binary data",
			req:  &pbrpc.DataViewRequest{Id: 3},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 3).Return(&models.DBUserData{
					UserID:        userID,
					Type:          "binary_data",
					Meta:          `{"content":"meta"}`,
					MinioObjectID: "obj",
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				s3.On("GetObject", mock.Anything, "obj").Return([]byte("encdata"), &minio.ObjectInfo{
					UserMetadata: map[string]string{"original-name": "file.txt"},
					ContentType:  "txt",
				}, nil)
				env.On("DecryptUserData", mock.Anything, mock.AnythingOfType("models.DBUserData"), []byte("mk")).
					Return([]byte("filedata"), nil)
			},
			wantErr: false,
		},
		{
			name: "permission denied",
			req:  &pbrpc.DataViewRequest{Id: 4},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 4).Return(&models.DBUserData{
					UserID: 99, // not userID
					Type:   "bank_card",
					Meta:   `{"content":"meta"}`,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "not found",
			req:  &pbrpc.DataViewRequest{Id: 5},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 5).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "decrypt error",
			req:  &pbrpc.DataViewRequest{Id: 6},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 6).Return(&models.DBUserData{
					UserID: userID,
					Type:   "bank_card",
					Meta:   `{"content":"meta"}`,
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("DecryptUserData", mock.Anything, mock.AnythingOfType("models.DBUserData"), []byte("mk")).
					Return(nil, errors.New("decrypt fail"))
			},
			wantErr: true,
		},
		{
			name: "unsupported type",
			req:  &pbrpc.DataViewRequest{Id: 7},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 7).Return(&models.DBUserData{
					UserID: userID,
					Type:   "unknown_type",
					Meta:   `{"content":"meta"}`,
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("DecryptUserData", mock.Anything, mock.AnythingOfType("models.DBUserData"), []byte("mk")).
					Return([]byte("data"), nil)
			},
			wantErr: true,
		},
		{
			name: "meta unmarshal error",
			req:  &pbrpc.DataViewRequest{Id: 8},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				st.On("GetUserData", mock.Anything, 8).Return(&models.DBUserData{
					UserID: userID,
					Type:   "bank_card",
					Meta:   `not_json`,
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("DecryptUserData", mock.Anything, mock.AnythingOfType("models.DBUserData"), []byte("mk")).
					Return([]byte("data"), nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := mocks.NewIStorage(t)
			s3 := mocks.NewS3Client(t)
			env := mocks.NewIEnvelope(t)
			km := mocks.NewKeyManagerInterface(t)

			tt.setupMocks(st, s3, env, km)
			srv := &ServerAdmin{
				Storage:    st,
				StorageS3:  s3,
				JWTConfig:  config.JWTConfig{},
				Envelope:   env,
				KeyManager: km,
			}
			_, err := srv.DataView(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
