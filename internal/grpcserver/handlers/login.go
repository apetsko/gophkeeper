package handlers

import (
	"context"
	"errors"
	"fmt"

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

	return &pbrpcu.LoginResponse{
		Token: token,
	}, nil
}
