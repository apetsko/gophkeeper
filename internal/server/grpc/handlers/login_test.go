package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/internal/mocks"
	"github.com/apetsko/gophkeeper/models"
	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
	"github.com/apetsko/gophkeeper/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServerAdmin_Login(t *testing.T) {
	const (
		userID   = 42
		username = "testuser"
		password = "password123"
		secret   = "secret"
	)
	// Generate a real hash for password check
	hash, _ := utils.HashPassword(password)

	tests := []struct {
		name       string
		req        *pbrpcu.LoginRequest
		setupMocks func(st *mocks.IStorage, km *mocks.KeyManagerInterface)
		wantErr    bool
	}{
		{
			name: "success",
			req:  &pbrpcu.LoginRequest{Username: username, Password: password},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {
				st.On("GetUser", mock.Anything, username).Return(&models.UserEntry{
					ID:           userID,
					Username:     username,
					PasswordHash: string(hash),
				}, nil)
				km.On("GetOrCreateMasterKey", mock.Anything, userID, password, mock.Anything).Return([]byte("mk"), nil)
			},
			wantErr: false,
		},
		{
			name:       "short username",
			req:        &pbrpcu.LoginRequest{Username: "ab", Password: password},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {},
			wantErr:    true,
		},
		{
			name:       "short password",
			req:        &pbrpcu.LoginRequest{Username: username, Password: "short"},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {},
			wantErr:    true,
		},
		{
			name: "user not found",
			req:  &pbrpcu.LoginRequest{Username: username, Password: password},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {
				st.On("GetUser", mock.Anything, username).Return(nil, models.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "wrong password",
			req:  &pbrpcu.LoginRequest{Username: username, Password: "wrongpass"},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {
				st.On("GetUser", mock.Anything, username).Return(&models.UserEntry{
					ID:           userID,
					Username:     username,
					PasswordHash: string(hash),
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "storage error",
			req:  &pbrpcu.LoginRequest{Username: username, Password: password},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {
				st.On("GetUser", mock.Anything, username).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "master key error",
			req:  &pbrpcu.LoginRequest{Username: username, Password: password},
			setupMocks: func(st *mocks.IStorage, km *mocks.KeyManagerInterface) {
				st.On("GetUser", mock.Anything, username).Return(&models.UserEntry{
					ID:           userID,
					Username:     username,
					PasswordHash: string(hash),
				}, nil)
				km.On("GetOrCreateMasterKey", mock.Anything, userID, password, mock.Anything).Return(nil, errors.New("mk error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := mocks.NewIStorage(t)
			km := mocks.NewKeyManagerInterface(t)
			tt.setupMocks(st, km)

			srv := &ServerAdmin{
				Storage:    st,
				JWTConfig:  config.JWTConfig{Secret: secret},
				KeyManager: km,
			}

			resp, err := srv.Login(context.Background(), tt.req)
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
