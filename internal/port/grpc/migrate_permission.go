package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) MigratePermissions(ctx context.Context, req *authPb.MigratePermissionsRequest) (*authPb.MigratePermissionsResponse, error) {

	err := s.svc.GetAuthService().MigratePermission(ctx, auth.MigratePermissionInput{
		Permissions: FromListPermissionProto(req.Permissions),
		Paths:       FromListPermissionPathProto(req.Paths),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.MigratePermissionsResponse{}, nil
}
