package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) UpdateRole(ctx context.Context, req *authPb.UpdateRoleRequest) (*authPb.UpdateRoleResponse, error) {

	_, err := s.svc.GetAuthService().UpdateRole(ctx, auth.UpdateRoleInput{
		ID:            req.Id,
		Name:          req.Name,
		PermissionIDs: req.PermissionIds,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.UpdateRoleResponse{}, nil
}
