package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) HasPermission(ctx context.Context, req *authPb.HasPermissionRequest) (*authPb.HasPermissionResponse, error) {

	res, err := s.svc.GetAuthService().HasPermission(ctx, auth.HasPermissionInput{
		UserID:   req.UserId,
		Resource: req.Resource,
		Action:   req.Action,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.HasPermissionResponse{
		Allowed: res,
	}, nil
}
