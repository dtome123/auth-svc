package grpc

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) UpdateRole(ctx context.Context, req *authPb.UpdateRoleRequest) (*authPb.UpdateRoleResponse, error) {
	var ras []domain.ResourceAction
	var permIDs []string
	for _, ra := range req.GetResourceActions() {
		ras = append(ras, domain.ResourceAction{
			PermissionID: ra.GetPermissionId(),
			Actions:      ra.GetActions(),
		})
		if ra.GetPermissionId() != "" {
			permIDs = append(permIDs, ra.GetPermissionId())
		}
	}

	_, err := s.svc.GetAuthService().UpdateRole(ctx, auth.UpdateRoleInput{
		ID:              req.GetId(),
		Name:            req.GetName(),
		ResourceActions: ras,
		PermissionIDs:   permIDs,
		Description:     req.GetDescription(),
		Scope:           req.GetScope(),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.UpdateRoleResponse{}, nil
}
