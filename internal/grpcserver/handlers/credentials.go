package handlers

import (
	"context"
	"fmt"

	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) Credentials(ctx context.Context, in *pbrpc.CredentialsRequest) (*pbrpc.CredentialsResponse, error) {
	fmt.Println(in)

	return &pbrpc.CredentialsResponse{
		Message: "ok",
	}, nil
}
