package handlers

import (
	"context"
	"fmt"

	pbrpcu "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
)

func (s *ServerAdmin) Login(ctx context.Context, in *pbrpcu.LoginRequest) (*pbrpcu.LoginResponse, error) {
	fmt.Println(in.Username)

	return &pbrpcu.LoginResponse{
		Token: "12345",
	}, nil
}
