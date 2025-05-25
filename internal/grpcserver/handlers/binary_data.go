package handlers

import (
	"context"
	"fmt"

	pbrpc "gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) BinaryData(ctx context.Context, in *pbrpc.BinaryDataRequest) (*pbrpc.BinaryDataResponse, error) {
	fmt.Println(in)

	return &pbrpc.BinaryDataResponse{
		Message: "ok",
	}, nil
}
