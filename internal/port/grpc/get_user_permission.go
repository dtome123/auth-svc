package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) GetUserPermissions(ctx context.Context, req *authPb.GetUserPermissionsRequest) (*authPb.GetUserPermissionsResponse, error) {

	res, err := s.svc.GetAuthService().GetUserPermissions(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &authPb.GetUserPermissionsResponse{
		Permissions: ToPermissionsProto(res),
	}, nil
}
