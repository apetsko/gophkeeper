package handlers

import (
	"context"
	"fmt"

	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) BankCard(ctx context.Context, in *pbrpc.BankCardRequest) (*pbrpc.BankCardResponse, error) {
	fmt.Println(in)

	return &pbrpc.BankCardResponse{
		Message: "ok",
	}, nil
}
