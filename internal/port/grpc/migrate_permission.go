package grpc

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) MigratePermissions(ctx context.Context, req *authPb.MigratePermissionsRequest) (*authPb.MigratePermissionsResponse, error) {
	var roles []*domain.Role
	for _, r := range req.GetRoles() {
		roles = append(roles, FromRoleProto(r))
	}

	err := s.svc.GetAuthService().MigratePermission(ctx, auth.MigratePermissionInput{
		Permissions: FromListPermissionProto(req.GetPermissions()),
		Roles:       roles,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.MigratePermissionsResponse{}, nil
}
