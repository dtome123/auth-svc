package grpc

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) CreateRole(ctx context.Context, req *authPb.CreateRoleRequest) (*authPb.CreateRoleResponse, error) {
	var ras []domain.ResourceAction
	for _, ra := range req.GetResourceActions() {
		ras = append(ras, domain.ResourceAction{
			PermissionID: ra.GetPermissionId(),
			Actions:      ra.GetActions(),
		})
	}
	id, err := s.svc.GetAuthService().CreateRole(ctx, auth.CreateRoleInput{
		Code:            req.Code,
		Name:            req.GetName(),
		ResourceActions: ras,
		Description:     req.GetDescription(),
		Scope:           req.GetScope(),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.CreateRoleResponse{
		Id: id,
	}, nil
}
