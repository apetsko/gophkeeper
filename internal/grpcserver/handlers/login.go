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

func (s *ServerAdmin) Login(ctx context.Context, in *pbrpcu.LoginRequest) (*pbrpcu.LoginResponse, error) {
	if len(in.Username) < 3 || len(in.Password) < 8 {
		return nil, fmt.Errorf("username and password must be at least 3 and 8 characters long")
	}

	user, err := s.storage.GetUser(ctx, in.Username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !password.CheckPasswordHash(in.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := jwt.GenerateJWT(user.ID, user.Username, s.jwtConfig.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	encryptedMK, err := s.keyManager.GetMasterKey(
		ctx,
		user.ID,
		in.Password,
		nil,
	)

	if err != nil {
		slog.Error("failed to generate encrypted master key: %w", err)

		return nil, errors.New("failed to generate encrypted master key")
	}

	// TODO: нужно записать в потокобезопасную мапу в памяти
	fmt.Println("encryptedMK: [", encryptedMK, "]")

	return &pbrpcu.LoginResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Token:    token,
	}, nil
}
