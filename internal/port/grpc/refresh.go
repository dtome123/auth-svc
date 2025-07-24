package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/gen/go/auth/v1"
)

func (GrpcServer) Refresh(ctx context.Context, req *authPb.RefreshRequest) (*authPb.RefreshResponse, error) {
	return nil, nil
}
		