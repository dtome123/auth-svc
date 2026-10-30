package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) GetRoleByCode(ctx context.Context, req *authPb.GetRoleByCodeRequest) (*authPb.GetRoleByCodeResponse, error) {
	role, err := s.svc.GetAuthService().GetRoleByCode(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}

	return &authPb.GetRoleByCodeResponse{
		Role: ToRoleProto(role),
	}, nil
}
