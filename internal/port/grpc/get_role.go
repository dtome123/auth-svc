package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) GetRole(ctx context.Context, req *authPb.GetRoleRequest) (*authPb.GetRoleResponse, error) {
	role, err := s.svc.GetAuthService().GetRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, err
	}

	return &authPb.GetRoleResponse{
		Role: ToRoleProto(role),
	}, nil
}
