package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) CreateRole(ctx context.Context, req *authPb.CreateRoleRequest) (*authPb.CreateRoleResponse, error) {

	id, err := s.svc.GetAuthService().CreateRole(ctx, auth.CreateRoleInput{
		Name:          req.Name,
		PermissionIDs: req.PermissionIds,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.CreateRoleResponse{
		Id: id,
	}, nil
}
