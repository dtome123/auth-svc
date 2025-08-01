package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) CheckPermissions(ctx context.Context, req *authPb.CheckPermissionsRequest) (*authPb.CheckPermissionsResponse, error) {

	res, err := s.svc.GetAuthService().CheckPermissions(ctx, auth.CheckPermissionsInput{
		UserID:          req.UserId,
		ActionResources: FromListActionResourceProto(req.Checks),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.CheckPermissionsResponse{
		Results: res,
	}, nil
}
