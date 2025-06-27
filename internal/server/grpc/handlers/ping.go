// Package handlers provides gRPC server handlers for managing user data operations,
// including creation, retrieval, update, and deletion of user records.
package handlers

import (
	"context"

	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

// Ping handles a health check request for the gRPC service.
//
// This method is used to verify that the server is alive and responding to requests.
// It does not perform any logic and always returns an empty successful response.
//
// Parameters:
// - ctx: The gRPC context (not used).
// - in: The PingRequest message (not used).
//
// Returns:
// - *pbrpc.PingResponse: An empty response indicating success.
// - error: Always nil.
func (s *ServerAdmin) Ping(_ context.Context, _ *pbrpc.PingRequest) (*pbrpc.PingResponse, error) {
	return &pbrpc.PingResponse{
		Message: "pong",
	}, nil
}
