package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/mocks"
	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServerAdmin_Signup(t *testing.T) {
	const (
		userID      = 42
		username    = "testuser"
		passwordStr = "password123"
		secret      = "secret"
	)
	tests := []struct {
		name         string
		req          *pbrpcu.SignupRequest
		setupStorage func(st *mocks.IStorage)
		setupKeys    func(km *mocks.KeyManagerInterface)
		wantErr      bool
	}{
		{
			name: "success",
			req:  &pbrpcu.SignupRequest{Username: username, Password: passwordStr},
			setupStorage: func(st *mocks.IStorage) {
				st.On("AddUser", mock.Anything, mock.AnythingOfType("*models.UserEntry")).Return(userID, nil)
			},
			setupKeys: func(km *mocks.KeyManagerInterface) {
				km.On("GetOrCreateMasterKey", mock.Anything, userID, passwordStr, mock.Anything).Return([]byte("key"), nil)
			},
			wantErr: false,
		},
		{
			name:         "short username",
			req:          &pbrpcu.SignupRequest{Username: "ab", Password: passwordStr},
			setupStorage: func(st *mocks.IStorage) {},
			setupKeys:    func(km *mocks.KeyManagerInterface) {},
			wantErr:      true,
		},
		{
			name:         "short password",
			req:          &pbrpcu.SignupRequest{Username: username, Password: "short"},
			setupStorage: func(st *mocks.IStorage) {},
			setupKeys:    func(km *mocks.KeyManagerInterface) {},
			wantErr:      true,
		},
		{
			name: "storage error",
			req:  &pbrpcu.SignupRequest{Username: username, Password: passwordStr},
			setupStorage: func(st *mocks.IStorage) {
				st.On("AddUser", mock.Anything, mock.AnythingOfType("*models.UserEntry")).Return(0, errors.New("db error"))
			},
			setupKeys: func(km *mocks.KeyManagerInterface) {},
			wantErr:   true,
		},
		{
			name: "key manager error",
			req:  &pbrpcu.SignupRequest{Username: username, Password: passwordStr},
			setupStorage: func(st *mocks.IStorage) {
				st.On("AddUser", mock.Anything, mock.AnythingOfType("*models.UserEntry")).Return(userID, nil)
			},
			setupKeys: func(km *mocks.KeyManagerInterface) {
				km.On("GetOrCreateMasterKey", mock.Anything, userID, passwordStr, mock.Anything).Return(nil, errors.New("key error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := mocks.NewIStorage(t)
			km := mocks.NewKeyManagerInterface(t)
			tt.setupStorage(st)
			tt.setupKeys(km)

			srv := &ServerAdmin{
				Storage:    st,
				JWTConfig:  config.JWTConfig{Secret: secret},
				KeyManager: km,
			}

			resp, err := srv.Signup(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, username, resp.Username)
				assert.NotEmpty(t, resp.Token)
			}
		})
	}
}
