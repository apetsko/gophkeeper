package handlers

import (
	"context"
	"fmt"

	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) DataDelete(ctx context.Context, in *pbrpc.DataDeleteRequest) (*pbrpc.DataDeleteResponse, error) {
	fmt.Println(in.GetId())

	return &pbrpc.DataDeleteResponse{
		Message: "ok",
	}, nil
}
