// Package handlers provides gRPC server handlers for managing user data operations,
// including creation, retrieval, update, and deletion of user records.
package handlers

import (
	"context"
	"fmt"

	"github.com/apetsko/gophkeeper/models"
	"github.com/apetsko/gophkeeper/pkg/jwt"
	"github.com/apetsko/gophkeeper/pkg/password"
	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
)

// Signup handles the gRPC request for user registration.
//
// This method validates the username and password, hashes the password,
// creates a new user record in the database, and generates a JWT token for the user.
//
// Parameters:
//   - ctx: The gRPC context.
//   - in: The SignupRequest message with user credentials.
//
// Returns:
//   - *pbrpcu.SignupResponse: User details and authentication token.
//   - error: An error if registration fails.
func (s *ServerAdmin) Signup(ctx context.Context, in *pbrpcu.SignupRequest) (*pbrpcu.SignupResponse, error) {
	if len(in.Username) < 3 || len(in.Password) < 8 {
		return nil, fmt.Errorf("username and password must be at least 3 and 8 characters long")
	}

	hash, err := password.HashPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.UserEntry{
		Username:     in.Username,
		PasswordHash: hash,
	}

	userID, err := s.Storage.AddUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := jwt.GenerateJWT(userID, in.Username, s.JWTConfig.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &pbrpcu.SignupResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Token:    token,
	}, nil
}
