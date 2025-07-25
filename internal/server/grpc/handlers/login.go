// Package handlers provides gRPC server handlers for managing user data operations,
// including creation, retrieval, update, and deletion of user records.
package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/apetsko/gophkeeper/models"
	"github.com/apetsko/gophkeeper/pkg/jwt"
	"github.com/apetsko/gophkeeper/pkg/password"
	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
)

// Login handles the gRPC request for user authentication.
//
// This method validates the username and password, checks credentials against the database,
// generates a JWT token upon successful authentication, and ensures the user's master key exists.
//
// Parameters:
// - ctx: The gRPC context.
// - in: The LoginRequest message with user credentials.
//
// Returns:
// - *pbrpcu.LoginResponse: User details and authentication token.
// - error: An error if authentication fails.
func (s *ServerAdmin) Login(ctx context.Context, in *pbrpcu.LoginRequest) (*pbrpcu.LoginResponse, error) {
	if len(in.Username) < 3 || len(in.Password) < 8 {
		return nil, fmt.Errorf("username and password must be at least 3 and 8 characters long")
	}

	user, err := s.Storage.GetUser(ctx, in.Username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !password.CheckPasswordHash(in.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := jwt.GenerateJWT(user.ID, user.Username, s.JWTConfig.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// TODO: нужно записать в потокобезопасную мапу в памяти
	_, errMasterKey := s.KeyManager.GetOrCreateMasterKey(
		ctx,
		user.ID,
		in.Password,
		nil,
	)

	if errMasterKey != nil {
		slog.Error("failed to generate encrypted master key: " + errMasterKey.Error())

		return nil, errors.New("failed to generate encrypted master key")
	}

	return &pbrpcu.LoginResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Token:    token,
	}, nil
}
