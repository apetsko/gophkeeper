package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/constants"
	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/models"
	pbc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/common"
	pbmodels "github.com/apetsko/gophkeeper/protogen/api/proto/v1/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServerAdmin_DataSave(t *testing.T) {
	const userID = 42
	ctx := context.WithValue(context.Background(), constants.UserID, userID)

	type testCase struct {
		name       string
		req        *pbrpc.DataSaveRequest
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
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Meta: &pbmodels.Meta{Content: "meta"},
				Data: &pbrpc.DataSaveRequest_BankCard{
					BankCard: &pbmodels.BankCard{
						CardNumber: "1234",
					},
				},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				st.On("SaveUserData", mock.Anything, mock.AnythingOfType("*models.DBUserData")).Return(1, nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(&models.EncryptedData{
					EncryptedData: []byte("enc"),
					DataNonce:     []byte("nonce"),
					EncryptedDek:  []byte("dek"),
					DekNonce:      []byte("dek_nonce"),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success credentials",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_CREDENTIALS,
				Meta: &pbmodels.Meta{Content: "meta"},
				Data: &pbrpc.DataSaveRequest_Credentials{
					Credentials: &pbmodels.Credentials{
						Login:    "l",
						Password: "p",
					},
				},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				st.On("SaveUserData", mock.Anything, mock.AnythingOfType("*models.DBUserData")).Return(1, nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(&models.EncryptedData{
					EncryptedData: []byte("enc"),
					DataNonce:     []byte("nonce"),
					EncryptedDek:  []byte("dek"),
					DekNonce:      []byte("dek_nonce"),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success binary data",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BINARY_DATA,
				Meta: &pbmodels.Meta{Content: "meta"},
				Data: &pbrpc.DataSaveRequest_BinaryData{
					BinaryData: &pbmodels.File{
						Name: "file.txt",
						Data: []byte("data"),
						Type: "txt",
					},
				},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(&models.EncryptedData{
					EncryptedData: []byte("enc"),
					DataNonce:     []byte("nonce"),
					EncryptedDek:  []byte("dek"),
					DekNonce:      []byte("dek_nonce"),
				}, nil)
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				s3.On("Upload", mock.Anything, mock.Anything, mock.AnythingOfType("*models.S3UploadData")).Return(nil, nil)
				st.On("SaveUserData", mock.Anything, mock.AnythingOfType("*models.DBUserData")).Return(1, nil)
			},
			wantErr: false,
		},
		{
			name: "missing user id",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {},
			wantErr:    true,
		},
		{
			name: "unsupported type",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_UNSPECIFIED,
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {},
			wantErr:    true,
		},
		{
			name: "key manager error",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Meta: &pbmodels.Meta{Content: "meta"},
				Data: &pbrpc.DataSaveRequest_BankCard{
					BankCard: &pbmodels.BankCard{
						CardNumber: "1234",
					},
				},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return(nil, errors.New("fail"))
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

			if tt.name == "missing user id" {
				ctx := context.Background()
				srv := &ServerAdmin{
					Storage:    st,
					StorageS3:  s3,
					JWTConfig:  config.JWTConfig{},
					Envelope:   nil,
					KeyManager: km,
				}
				_, err := srv.DataSave(ctx, tt.req)
				assert.Error(t, err)
				return
			}

			tt.setupMocks(st, s3, env, km)
			srv := &ServerAdmin{
				Storage:    st,
				StorageS3:  s3,
				JWTConfig:  config.JWTConfig{},
				Envelope:   env,
				KeyManager: km,
			}
			_, err := srv.DataSave(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServerAdmin_DataSave_ErrorBranches(t *testing.T) {
	const userID = 42
	ctx := context.WithValue(context.Background(), constants.UserID, userID)

	cases := []struct {
		name       string
		req        *pbrpc.DataSaveRequest
		setupMocks func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface)
		wantErr    string
	}{
		{
			name: "bank card nil",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Data: &pbrpc.DataSaveRequest_BankCard{BankCard: nil},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
			},
			wantErr: "отсутствуют данные банковской карты",
		},
		{
			name: "credentials nil",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_CREDENTIALS,
				Data: &pbrpc.DataSaveRequest_Credentials{Credentials: nil},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
			},
			wantErr: "отсутствуют учетные данные",
		},
		{
			name: "file nil",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BINARY_DATA,
				Data: &pbrpc.DataSaveRequest_BinaryData{BinaryData: nil},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
			},
			wantErr: "отсутствуют данные файла",
		},
		{
			name: "envelope encrypt error (binary)",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BINARY_DATA,
				Data: &pbrpc.DataSaveRequest_BinaryData{BinaryData: &pbmodels.File{Name: "f", Data: []byte("d")}},
				Meta: &pbmodels.Meta{},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
			},
			wantErr: "failed to encrypt binary data",
		},
		{
			name: "s3 upload error",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BINARY_DATA,
				Data: &pbrpc.DataSaveRequest_BinaryData{BinaryData: &pbmodels.File{Name: "f", Data: []byte("d")}},
				Meta: &pbmodels.Meta{},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(&models.EncryptedData{
					EncryptedData: []byte("enc"),
					DataNonce:     []byte("nonce"),
					EncryptedDek:  []byte("dek"),
					DekNonce:      []byte("dek_nonce"),
				}, nil)
				s3.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
			},
			wantErr: "failed to upload file to MinIO",
		},
		{
			name: "save user data error (binary)",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BINARY_DATA,
				Data: &pbrpc.DataSaveRequest_BinaryData{BinaryData: &pbmodels.File{Name: "f", Data: []byte("d")}},
				Meta: &pbmodels.Meta{},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(&models.EncryptedData{
					EncryptedData: []byte("enc"),
					DataNonce:     []byte("nonce"),
					EncryptedDek:  []byte("dek"),
					DekNonce:      []byte("dek_nonce"),
				}, nil)
				s3.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				st.On("SaveUserData", mock.Anything, mock.Anything).Return(0, errors.New("fail"))
			},
			wantErr: "fail",
		},
		{
			name: "envelope encrypt error (non-binary)",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Data: &pbrpc.DataSaveRequest_BankCard{BankCard: &pbmodels.BankCard{CardNumber: "1234"}},
				Meta: &pbmodels.Meta{},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
			},
			wantErr: "encrypt error",
		},
		{
			name: "save user data error (non-binary)",
			req: &pbrpc.DataSaveRequest{
				Type: pbc.DataType_DATA_TYPE_BANK_CARD,
				Data: &pbrpc.DataSaveRequest_BankCard{BankCard: &pbmodels.BankCard{CardNumber: "1234"}},
				Meta: &pbmodels.Meta{},
			},
			setupMocks: func(st *mocks.IStorage, s3 *mocks.S3Client, env *mocks.IEnvelope, km *mocks.KeyManagerInterface) {
				km.On("GetMasterKey", mock.Anything, userID).Return([]byte("mk"), nil)
				env.On("EncryptUserData", mock.Anything, mock.Anything, mock.Anything).Return(&models.EncryptedData{
					EncryptedData: []byte("enc"),
					DataNonce:     []byte("nonce"),
					EncryptedDek:  []byte("dek"),
					DekNonce:      []byte("dek_nonce"),
				}, nil)
				st.On("SaveUserData", mock.Anything, mock.Anything).Return(0, errors.New("fail"))
			},
			wantErr: "fail",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			st := mocks.NewIStorage(t)
			s3 := mocks.NewS3Client(t)
			env := mocks.NewIEnvelope(t)
			km := mocks.NewKeyManagerInterface(t)
			if tt.setupMocks != nil {
				tt.setupMocks(st, s3, env, km)
			}
			srv := &ServerAdmin{
				Storage:    st,
				StorageS3:  s3,
				JWTConfig:  config.JWTConfig{},
				Envelope:   env,
				KeyManager: km,
			}
			_, err := srv.DataSave(ctx, tt.req)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// brokenProto is a proto.Message that always fails to marshal
type brokenProto struct{ pbmodels.BankCard }

func (b *brokenProto) Reset()                   {}
func (b *brokenProto) String() string           { return "" }
func (b *brokenProto) ProtoMessage()            {}
func (b *brokenProto) Marshal() ([]byte, error) { return nil, errors.New("fail") }
