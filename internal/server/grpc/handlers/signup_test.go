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
		name       string
		req        *pbrpcu.SignupRequest
		setupMocks func(st *mocks.IStorage)
		wantErr    bool
	}{
		{
			name: "success",
			req:  &pbrpcu.SignupRequest{Username: username, Password: passwordStr},
			setupMocks: func(st *mocks.IStorage) {
				st.On("AddUser", mock.Anything, mock.AnythingOfType("*models.UserEntry")).Return(userID, nil)
			},
			wantErr: false,
		},
		{
			name:       "short username",
			req:        &pbrpcu.SignupRequest{Username: "ab", Password: passwordStr},
			setupMocks: func(st *mocks.IStorage) {},
			wantErr:    true,
		},
		{
			name:       "short password",
			req:        &pbrpcu.SignupRequest{Username: username, Password: "short"},
			setupMocks: func(st *mocks.IStorage) {},
			wantErr:    true,
		},
		{
			name: "storage error",
			req:  &pbrpcu.SignupRequest{Username: username, Password: passwordStr},
			setupMocks: func(st *mocks.IStorage) {
				st.On("AddUser", mock.Anything, mock.AnythingOfType("*models.UserEntry")).Return(0, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := mocks.NewIStorage(t)
			tt.setupMocks(st)

			srv := &ServerAdmin{
				Storage:   st,
				JWTConfig: config.JWTConfig{Secret: secret},
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
