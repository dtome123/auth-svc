package grpc

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) UpdateRole(ctx context.Context, req *authPb.UpdateRoleRequest) (*authPb.UpdateRoleResponse, error) {
	var ras []domain.ResourceAction
	for _, ra := range req.GetResourceActions() {
		ras = append(ras, domain.ResourceAction{
			PermissionID: ra.GetPermissionId(),
			Actions:      ra.GetActions(),
		})
	}

	_, err := s.svc.GetAuthService().UpdateRole(ctx, auth.UpdateRoleInput{
		ID:              req.GetId(),
		Name:            req.GetName(),
		ResourceActions: ras,
		Description:     req.GetDescription(),
		Scope:           req.GetScope(),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.UpdateRoleResponse{}, nil
}
