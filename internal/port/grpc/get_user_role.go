package grpc

import (
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) GetUserRoles(ctx context.Context, req *authPb.GetUserRolesRequest) (*authPb.GetUserRolesResponse, error) {

	res, err := s.svc.GetAuthService().GetUserRoles(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &authPb.GetUserRolesResponse{
		Roles: ToListRoleProto(res),
	}, nil
}
