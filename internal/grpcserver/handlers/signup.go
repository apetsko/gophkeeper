package handlers

import (
	"context"
	"fmt"

	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
)

func (s *ServerAdmin) Signup(ctx context.Context, in *pbrpcu.SignupRequest) (*pbrpcu.SignupResponse, error) {
	fmt.Println(in.Username)

	return &pbrpcu.SignupResponse{
		Token: "12345",
	}, nil
}
