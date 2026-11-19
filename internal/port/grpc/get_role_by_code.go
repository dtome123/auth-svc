package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) GetRoleByCode(ctx context.Context, req *authPb.GetRoleByCodeRequest) (*authPb.GetRoleByCodeResponse, error) {

	scope := make(map[string][]string)
	for k, v := range req.GetScope() {
		if v != nil {
			scope[k] = v.GetValues()
		}
	}

	role, err := s.svc.GetAuthService().GetRoleByCode(ctx, req.GetCode(), scope)
	if err != nil {
		return nil, err
	}

	return &authPb.GetRoleByCodeResponse{
		Role: ToRoleProto(role),
	}, nil
}
