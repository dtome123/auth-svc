package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/gen/go/auth/v1"
)

func (GrpcServer) Sign(ctx context.Context, req *authPb.SignRequest) (*authPb.SignResponse, error) {
	return nil, nil
}
